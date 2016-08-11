package abstract

import(
	"fmt"
	"io"
	"sync"
)

/*
 *  output := NewOuput(os.Stdin, os.Stderr)
 *  output.Write("file size is %s", stat.Size())
 *  output.Done; output.Wait()
*/

// write in a separate goroutine to output writers (stdout, stderr)
type Output struct {
	outwrite	io.Writer
	ewrite		io.Writer

	out		chan *Message
	eout		chan *Message

	waitGroup	sync.WaitGroup
}

func NewOutput(out, e io.Writer) (output *Output) {
	output = &Output { outwrite: out, ewrite: e}
	output.out = make(chan *Message, GOROUTINE_LIMIT_RESOURCE)
	output.eout = make(chan *Message, GOROUTINE_LIMIT_RESOURCE)

	output.waitGroup.Add(1)
	go func() {
	infinite:
		for {
			select {
			case message, ok := <-output.out:
				if !ok { break infinite }
				fmt.Fprintf(output.outwrite, message.Format, message.Values...)
			case emessage, ok := <-output.eout:
				if !ok { break infinite }
				fmt.Fprintf(output.ewrite, emessage.Format, emessage.Values...)
			}
		}
		output.waitGroup.Done()
	}()

	return
}

func (output *Output) Write(format string, values ...interface{}) {
	output.out <-&Message{ Format: format, Values: values }
}

func (output *Output) WriteE(e error) bool {
	if e == nil { return false }
	output.eout <-&Message{ Format: "%v\n", Values: []interface{} { e }  }
	return true
}

func (output *Output) WriteEMessage(e error, format string, values ...interface{}) bool {
	if e == nil { return false }
	format = "%v: " + format
	v := make([]interface{}, len(values)); v = append(v, e); v = append(v, values...)
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

func (output *Output) Wait() {
	output.waitGroup.Wait()
}

