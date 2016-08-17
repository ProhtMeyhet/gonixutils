package abstract

import(
	"bytes"
	"io"
	"sort"
	"strconv"
	"text/tabwriter"
)

// buffers the output, but not e, until Done(), sorts that and outputs it to the given output
type SortedOutput struct {
	Output

	buffer		map[string]*Message
	appendBuffer	map[string]*Message
	intsortkey	int
	lastSortKey	string

	usekey		int
	sortTransformator func(string) string
}

func NewSortedOutput(out, e io.Writer) (output OutputInterface) {
	output = &SortedOutput{}
	output.Initialise(out, e)
	return
}

func NewSortedTabbedOutput(out, e io.Writer) OutputInterface {
	output := &SortedOutput{ }
	output.tabwriter = tabwriter.NewWriter(out, 32, 0, 0, ' ', 0)
	output.terminalInfo = new(TerminalInfo)
	output.delimiter = "\t"
	output.Initialise(output.tabwriter, e)
	return output
}

// TODO add a copy function
func (output *SortedOutput) NewSubBuffer(name string, key int) OutputInterface {
	buffer := &SortedOutput{}
	buffer.linesManual = output.linesManual
	buffer.printSubBufferNames = output.printSubBufferNames
	buffer.orderBySubBufferNames = output.orderBySubBufferNames
	buffer.orderBySubBufferKeys = output.orderBySubBufferKeys
	buffer.sortReversed = output.sortReversed
	buffer.sortTransformator = output.sortTransformator

	// TODO is this correct? or is there a better alternative?
	bufferbyte := make([]byte, 0);	bytesBuffer := bytes.NewBuffer(bufferbyte)

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

func (output *SortedOutput) Initialise(out, e io.Writer) {
	output.buffer = make(map[string]*Message)
	output.appendBuffer = make(map[string]*Message)
	output.Output.Initialise(out, e)
}

// soft sort. tries to use the first value as stringkey
func (output *SortedOutput) Write(format string, values ...interface{}) {
	if len(values) == 0 || len(values) - 1 < output.usekey { output.Output.Write(format, values...); return }
	sortkey, ok := values[output.usekey].(string); if !ok { output.intsortkey++; sortkey = strconv.Itoa(output.intsortkey) }
	output.WriteSorted(format, sortkey, values...)
}

// append to the last sort key
func (output *SortedOutput) Append(format string, values ...interface{}) {
	output.appendBuffer[output.lastSortKey] = &Message{ Format: format, Values: values }
}

func (output *SortedOutput) WriteSorted(format, sortkey string, values ...interface{}) {
	output.lastSortKey = sortkey
	output.innerMutex.Lock()
	output.buffer[sortkey] = &Message{ Format: format, Values: values }
	output.innerMutex.Unlock()
}

func (output *SortedOutput) Done() {
	if output.done { return }

	sorted := make([]string, len(output.buffer)); i := 0

	// adding and removing while ranging is undefined behaviour (sic!)
	add := make(map[string]*Message, 0); remove := make([]string, 0)

	for key, _ := range output.buffer {
		if output.sortTransformator != nil {
			newKey := output.sortTransformator(key)
			sorted[i] = newKey
			if key != newKey {
				oldvalue := output.buffer[key]; remove = append(remove, key); add[sorted[i]] = oldvalue
			}
		} else {
			sorted[i] = key
		}
		i++
	}

	for key, value := range add {
		output.buffer[key] = value
	}

	for _, value := range remove {
		delete(output.buffer, value)
	}

	if output.sortReversed {
		sort.Sort(sort.Reverse(sort.StringSlice(sorted)))
	} else {
		sort.Strings(sorted)
	}

	for _, key := range sorted {
		output.Output.Write(output.buffer[key].Format, output.buffer[key].Values...)
	}

	for _, value := range output.appendBuffer {
		output.Output.Write(value.Format, value.Values...)
	}

	output.Output.Done()
}

func (output *SortedOutput) SortKey() int {
	return output.usekey
}

func (output *SortedOutput) SetSortKey(to int) {
	output.usekey = to
}

func (output *SortedOutput) SetSortTransformator(to func(string) string) {
	output.sortTransformator = to
}
