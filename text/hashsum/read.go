package hash

import(
	"io"

	"github.com/ProhtMeyhet/libgosimpleton/simpleton"
	"github.com/ProhtMeyhet/libgosimpleton/iotool"
)

func ReadFile(helper *iotool.FileHelper, buffers chan NamedBuffer, path string) (e error) {
	handlerBuffer := make([]byte, helper.ReadSize())
	return readFile(helper, buffers, handlerBuffer, path, true)
}

// reads file sequential; errors are reported via helper.E
func ReadFiles(helper *iotool.FileHelper, buffers chan NamedBuffer, paths <-chan string) {
	handlerBuffer := make([]byte, helper.ReadSize()); first := true
	for path := range paths {
		if e := readFile(helper, buffers, handlerBuffer, path, first); e != nil {
			// has already been raised
			if iotool.IsNotExist(e) { continue }
			helper.RaiseError(path, e)
		}
		first = false
	}
}

// reads file sequential; errors are reported via helper.E
func ReadFilesFromList(helper *iotool.FileHelper, buffers chan NamedBuffer, paths ...string) {
	ReadFiles(helper, buffers, simpleton.StringListToChannel(paths...))
}

// to avoid allocation in ReadFiles
func readFile(helper *iotool.FileHelper, buffers chan NamedBuffer, handlerBuffer []byte, path string, first bool) (e error) {
	handler, e := iotool.Open(helper, path); if e != nil { return }; defer handler.Close()
	namedBuffer := NewNamedBuffer(path); namedBuffer.buffer = make([]byte, len(handlerBuffer))
	if !first { namedBuffer.reset = true }; read := 0 // avoid e shadowing

infinite:
	for {
		read, e = handler.Read(handlerBuffer)
		if e != nil {
			if e == io.EOF { e = nil; break infinite }
			// TODO find out what errors can happen here and handle them
			break infinite
		}

	// FIXME i thought the following line would copy. doesn't seem to be or bug.
	//	buffer <-handlerBuffer[:read]
		namedBuffer.read = read
		copy(namedBuffer.buffer, handlerBuffer[:read])
		buffers <-namedBuffer
		namedBuffer.reset = false
	}

	namedBuffer.done = true
	buffers <-namedBuffer

	return
}
