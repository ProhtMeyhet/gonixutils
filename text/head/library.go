package head

import(
	"io"

//	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/iotool/ioreaders"
//	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
	"github.com/ProhtMeyhet/gonixutils/text/cat"
)

// unix head
func Head(input *Input) (exitCode uint8) {
	limit := func(reader io.Reader) io.Reader {
		return ioreaders.NewLimitLinesReader(reader, input.Max)
	}

	if input.Bytes {
		limit = func(reader io.Reader) io.Reader {
			return ioreaders.NewLimitBytesReader(reader, input.Max)
		}
	} else if input.Runes {
		limit = func(reader io.Reader) io.Reader {
			return ioreaders.NewLimitRunesReader(reader, input.Max)
		}
	}

	output := abstract.NewOutput(input.Stdout, input.Stderr)
	if len(input.Paths) > 1 { output.TogglePrintSubBufferNames() }
	helper := prepareFileHelper(input, output, &exitCode)

	e := cat.CopyFilesFilteredTo(output, helper, limit, input.Paths...); if e != nil && exitCode == 0 {
		// TODO parse the error and set exitCode accordingly
		exitCode = abstract.ERROR_UNHANDLED
	}


	output.Done(); output.Wait(); return
}
