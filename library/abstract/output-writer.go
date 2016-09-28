package abstract

import(
	"bytes"
	"fmt"
	"io"
	"strings"
//	"sort"
	"sync"
	"text/tabwriter"
)

/*
 *  output := NewOuput(os.Stdin, os.Stderr)
 *  output.Write("file size is %s", stat.Size())
 *  output.Done; output.Wait()
*/

// UNIX was not designed to stop you from doing stupid things, because that would also stop you from doing clever things.
//
//   — Doug Gwyn

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
	flushing		bool

	newLinePrintedOnce	bool
	delimiter		string

	// do not add line breaks
	linesManual		bool

	// subBuffer stuff
	subBuffers		[]OutputInterface
	subBufferNames		[]string
	subBufferKeys		[]int
	printSubBufferNames	bool
	orderBySubBufferNames	bool
	orderBySubBufferKeys	bool

	sortReversed		bool

	onceDo			sync.Once
	onceWait		sync.Once

	// called once in done
	done			bool
	waited			bool

	// public Lock() for public, inner mutex to guard assignment to maps
	innerMutex		sync.Mutex
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

// TODO add a possibilty to asynchronously flushing the buffer to it's parent when done
// TODO add a copy function
// return a new subbuffer that is automatically flushed to the main buffer when done
// name can be printed, key can be used for keeping a sort order in output
func (output *Output) NewSubBuffer(name string, key int) OutputInterface {
	buffer := &Output{}
	buffer.linesManual = output.linesManual
	buffer.printSubBufferNames = output.printSubBufferNames
	buffer.orderBySubBufferNames = output.orderBySubBufferNames
	buffer.orderBySubBufferKeys = output.orderBySubBufferKeys
	buffer.sortReversed = output.sortReversed

	// TODO is this correct? or is there a better alternative?
	bufferbyte := make([]byte, 0); bytesBuffer := bytes.NewBuffer(bufferbyte)

	if output.tabwriter != nil {
		buffer.tabwriter = tabwriter.NewWriter(bytesBuffer, 32, 0, 0, ' ', 0)
		buffer.terminalInfo = new(TerminalInfo)
		buffer.delimiter = "\t"
	}

	buffer.Initialise(bytesBuffer, output.ewrite)

	output.innerMutex.Lock()
	output.subBuffers = append(output.subBuffers, buffer)
	output.subBufferNames = append(output.subBufferNames, name)
	output.subBufferKeys = append(output.subBufferKeys, key)
	output.innerMutex.Unlock()

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

// write through
func (output *Output) Write(values []byte) (int, error) {
	return output.outwrite.Write(values)
}

func (output *Output) WriteFormatted(format string, values ...interface{}) {
	output.out <-&Message{ Format: format, Values: values }
}

// write unsorted
func (output *Output) WriteSorted(format, sortkey string, values ...interface{}) {
	output.WriteFormatted(format, values...)
}

func (output *Output) Append(format string, values ...interface{}) {
	output.WriteFormatted(format, values...)
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
	output.onceDo.Do(func() {
		close(output.out)
		close(output.eout)
		output.done = true
	})
}

func (output *Output) tabWrite(toWrite string) {
/*	// FIXME: please write a real tabwriter, instead of trying to fix this one
	// FIXME: try to test if the last one contains delimiter and print \n
	if strings.Contains(toWrite, output.delimiter) &&
		output.writeCount + 2 * output.tabBufferMax > output.terminalInfo.Width() -10 {
		toWrite += "\n"
	}*/

	//output.id++; toWrite = fmt.Sprintf("%v %v %v %v! ", output.id, output.writeCount, columnLen, output.terminalInfo.Width())
	// FIXME output.terminalInfo.Width() - 20 seems to be a good value...
	if output.writeCount + output.tabBufferMax > output.terminalInfo.Width() - 20 {
		toWrite += "\n"; output.writeCount = 0; output.newLinePrintedOnce = true
	} else { output.writeCount += output.tabBufferMax }
	fmt.Fprint(output.tabwriter, toWrite)
}

func (output *Output) Wait() {
	// lock is here for synchronisation
	output.Lock()
	output.onceWait.Do(func() {
		output.waitGroup.Wait()
	})

	if output.waited { output.Unlock(); return }

	// must be done here to ensure no one writes any more
	if output.tabwriter != nil {
		output.flushing = true
		// FIXME : tabwriter.Init() results in endless loop!
		if output.linesManual { output.tabBufferMax = 5 }
		output.tabwriter = tabwriter.NewWriter(output.outwrite, output.tabBufferMax, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
//		output.tabwriter = tabwriter.NewWriter(output.outwrite, output.tabBufferMax, 0, 1, '-', tabwriter.DiscardEmptyColumns|tabwriter.Debug)
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
/* TODO 
	if output.orderBySubBufferNames {
	    for key, name := range output.subBufferNames { fmt.Printf("%v %v\n", key, name) }
		sort.Strings(output.subBufferNames); i := 0
	    for key, name := range output.subBufferNames { fmt.Printf("%v %v\n", key, name) }
		for key, name := range output.subBufferNames {
			sub := output.subBuffers[key]
			if output.printSubBufferNames { fmt.Fprint(output.outwrite, name) }
			sub.Done(); sub.Wait()
			sub.writeTo(output.outwrite)
			if output.printSubBufferNames && i+1 != len(output.subBuffers){ fmt.Fprintln(output.outwrite) }
			i++
		}
	} else if output.orderBySubBufferKeys {
		sort.Ints(output.subBufferKeys); i := 0
		for key, _ := range output.subBufferKeys {
			sub := output.subBuffers[key]
			if output.printSubBufferNames { fmt.Fprint(output.outwrite, output.subBufferNames[key]) }
			sub.Done(); sub.Wait()
			sub.writeTo(output.outwrite)
			if output.printSubBufferNames && i+1 != len(output.subBuffers){ fmt.Fprintln(output.outwrite) }
			i++
		}
	} else { */
		for i := 0; i < len(output.subBuffers); i++ {
			sub := output.subBuffers[i]
			sub.Lock()
			if output.printSubBufferNames { fmt.Fprint(output.outwrite, output.subBufferNames[i]) }
			sub.Unlock(); sub.Done(); sub.Wait(); sub.Lock()
			sub.writeTo(output.outwrite)
			if output.printSubBufferNames && i+1 != len(output.subBuffers){ fmt.Fprintln(output.outwrite) }
			sub.Unlock()
		}
//	}

	output.waited = true
	output.Unlock()
}

func (output *Output) writeTo(to io.Writer) {
	if writeTo, ok := output.outwrite.(io.WriterTo); ok {
		writeTo.WriteTo(to)
	}
}

func (output *Output) Reset() {
	output.innerMutex.Lock(); output.Lock()
	output.onceDo = sync.Once{}
	output.onceWait = sync.Once{}
	output.subBuffers = make([]OutputInterface, 0)
	output.subBufferNames = make([]string, 0)
	output.subBufferKeys = make([]int, 0)
	output.tabBuffer = make([]string, 0)
	output.done = false
	output.waited = false
	output.Unlock(); output.innerMutex.Unlock()
}

func (output *Output) SortReversed() bool { return output.sortReversed }
func (output *Output) ToggleSortReversed() { output.sortReversed = !output.sortReversed }
func (output *Output) PrintSubBufferNames() bool { return output.printSubBufferNames }
func (output *Output) TogglePrintSubBufferNames() { output.printSubBufferNames = !output.printSubBufferNames }
func (output *Output) OrderBySubBufferNames() bool { return output.orderBySubBufferNames }
func (output *Output) ToggleOrderBySubBufferNames() { output.orderBySubBufferNames = !output.orderBySubBufferNames }
func (output *Output) OrderBySubBufferKeys() bool { return output.orderBySubBufferKeys }
func (output *Output) ToggleOrderBySubBufferKeys() { output.orderBySubBufferKeys = !output.orderBySubBufferKeys }

// output is a reduced version of sorted output
func (output *Output) SortKey() int { return -1 }
func (output *Output) SetSortKey(to int) { }
func (output *Output) SetSortTransformator(to func(string) string) { }

func (output *Output) LinesManual() bool { return output.linesManual }
func (output *Output) ToggleLinesManual() { output.linesManual = !output.linesManual }
