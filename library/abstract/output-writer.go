package abstract

import(
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"text/tabwriter"
	"os"
)

/*
 *  output := NewOuput(os.Stdin, os.Stderr)
 *  output.Write("file size is %s", stat.Size())
 *  output.Done; output.Wait()
*/

// write in a separate goroutine to output writers (stdout, stderr)
type Output struct {
	sync.Mutex

	outwrite	io.Writer
	ewrite		io.Writer

	out		chan *Message
	eout		chan *Message

	terminalInfo	*TerminalInfo
	tabwriter	*tabwriter.Writer
	tabBuffer	[]string
	// maximum tabs
	tabMax		int
	tabBufferMax	int
	writeCount	int

	waitGroup	sync.WaitGroup

	// currently flushing
	flushing	bool

	newLinePrintedOnce	bool
	delimiter		string

	// do not add line breaks
	linesManual		bool

	subBuffers	[]OutputInterface

	// public Lock() for public, inner mutex to guard assignment to maps
	innerMutex	sync.Mutex
}

func NewOutput(out, e io.Writer) (output OutputInterface) {
	output = &Output { }
	output.Initialise(out, e)
	return
}

func NewTabbedOutput(out, e io.Writer) OutputInterface {
	output := &Output{ }
	output.tabwriter = tabwriter.NewWriter(out, 32, 0, 0, ' ', 0)
	output.terminalInfo = new(TerminalInfo)
//	output.tabMax = output.terminalInfo.Width() / (7 + 3)
	output.Initialise(output.tabwriter, e)
	output.delimiter = "\t"
	return output
}

// TODO add a copy function
// return a new subbuffer that is automatically flushed to the main buffer when done
func (output *Output) NewSubBuffer() OutputInterface {
	buffer := &Output{}
	bufferbyte := make([]byte, 0) // TODO is this correct? or is there a better alternative?
	bytesBuffer := bytes.NewBuffer(bufferbyte)
	if output.tabwriter != nil {
		buffer.tabwriter = tabwriter.NewWriter(bytesBuffer, 32, 0, 0, ' ', 0)
		buffer.terminalInfo = new(TerminalInfo)
		buffer.delimiter = "\t"
	}
	buffer.Initialise(bytesBuffer, output.ewrite)
	output.subBuffers = append(output.subBuffers, buffer)
	buffer.linesManual = output.linesManual
	return buffer
}

func (output *Output) Initialise(out, e io.Writer) {
	output.outwrite = out; output.ewrite = e
	output.out = make(chan *Message, 10); output.eout = make(chan *Message, 2)

	output.waitGroup.Add(1)
	go func() {
	infinite:
		for {
			select {
			case message, ok := <-output.out:
				if !ok { break infinite }
			// FIXME write a generic wrapper around text/tabwriter that just takes a length argument
			// describing the max line (eg terminal) length and from there on writes columns correctly
				if output.tabwriter != nil && !output.flushing {
					pleaseWrite := fmt.Sprintf(message.Format, message.Values...) + output.delimiter
					if output.linesManual && strings.HasSuffix(pleaseWrite, "\n" + output.delimiter) {
						pleaseWrite = pleaseWrite[:len(pleaseWrite)-2] + "\n"
					}
					output.tabBuffer = append(output.tabBuffer, pleaseWrite)
					writeLen := len(pleaseWrite); if strings.HasPrefix(pleaseWrite, TERMINAL_COLOR_JUNK) {
						writeLen -= len(TERMINAL_COLOR_JUNK_LENGTH) + len(TERMINAL_COLOR_RESET)
					}
					if writeLen > output.tabBufferMax { output.tabBufferMax = len(pleaseWrite) }
				} else {
					fmt.Fprintf(output.outwrite, message.Format, message.Values...)
				}
			case emessage, ok := <-output.eout:
				if ok {
					fmt.Fprintf(output.ewrite, emessage.Format, emessage.Values...)
				}
			}
		}

		// read out all errors
	infiniteE:
		for {
			select {
			case emessage, ok := <-output.eout:
				if !ok { break infiniteE }
				fmt.Fprintf(output.ewrite, emessage.Format, emessage.Values...)
			}
		}
		output.waitGroup.Done()
	}()
}

func (output *Output) Write(format string, values ...interface{}) {
	output.out <-&Message{ Format: format, Values: values }
}

// write unsorted
func (output *Output) WriteSorted(format, sortkey string, values ...interface{}) {
	output.Write(format, values...)
}

func (output *Output) Append(format string, values ...interface{}) {
	output.Write(format, values...)
}

func (output *Output) WriteE(e error) bool {
	if e == nil { return false }
	output.eout <-&Message{ Format: "%v\n", Values: []interface{} { e }  }
	return true
}

func (output *Output) WriteEMessage(e error, format string, values ...interface{}) bool {
	if e == nil { return false }
	format = "%v: " + format + "\n"
	v := make([]interface{}, 0); v = append(v, e); v = append(v, values...)
	output.eout <-&Message{ Format: format, Values: v }
	return true
}

func (output *Output) WriteError(format string, values ...interface{}) {
	output.eout <-&Message{ Format: format, Values: values }
}

func (output *Output) Done() {
	close(output.out)
	close(output.eout)
}

func (output *Output) tabWrite(toWrite string) {
	//output.id++; toWrite = fmt.Sprintf("%v %v %v %v! ", output.id, output.writeCount, columnLen, output.terminalInfo.Width())
	if output.writeCount + output.tabBufferMax > output.terminalInfo.Width() {
		toWrite += "\n"; output.writeCount = 0; output.newLinePrintedOnce = true
	} else { output.writeCount += output.tabBufferMax }
	fmt.Fprint(output.tabwriter, toWrite)
}

func (output *Output) Wait() {
	output.waitGroup.Wait()

	// must be done here to ensure no one writes any more
	if output.tabwriter != nil {
		output.flushing = true
		// FIXME : tabwriter.Init() results in endless loop!
		if output.linesManual { output.tabBufferMax = 5 }
		output.tabwriter = tabwriter.NewWriter(output.outwrite, output.tabBufferMax, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
		for _, value := range output.tabBuffer {
			if !output.linesManual {
				output.newLinePrintedOnce = false
				output.tabWrite(value)
			} else {
				output.newLinePrintedOnce = true
				fmt.Fprint(output.tabwriter, value)
			}
		}; output.tabBuffer = make([]string, 0)
		// if line wasn't long enough and something was actually written
		if !output.newLinePrintedOnce && output.writeCount > 0 { fmt.Fprintf(output.tabwriter, "\n") }
		output.tabwriter.Flush()
	}

	for _, sub := range output.subBuffers {
		sub.Done(); sub.Wait()
		sub.writeTo(output.outwrite)
	}
}

func (output *Output) writeTo(to io.Writer) {
	if writeTo, ok := output.outwrite.(io.WriterTo); ok {
		writeTo.WriteTo(os.Stdout)
	}
}

// output is a reduced version of sorted output
func (output *Output) SortKey() int { return -1 }
func (output *Output) SetSortKey(to int) { }
func (output *Output) SetSortTransformator(to func(string) string) { }

func (output *Output) LinesManual() bool { return output.linesManual }
func (output *Output) ToggleLinesManual() { output.linesManual = !output.linesManual }
