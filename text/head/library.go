package head

import(
	"io"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/iotool/ioreaders"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
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
	helper := prepareFileHelper(input, output, &exitCode)

	exitCode = Limit(input, output, helper, limit)

	output.Done(); output.Wait()
	return
}

func Limit(input *Input, output abstract.OutputInterface, helper *iotool.FileHelper,
		limit func(io.Reader) io.Reader) (exitCode uint8) {
	key := 0
	parallel.ReadFilesFilteredSequential(helper, input.Paths, limit, func(buffers chan *iotool.NamedBuffer) {
		i := 0
		if !input.Quiet && len(input.Paths) > 1 && key == 0 { output.Write("==>%v<==\n", input.Paths[0]) }
		for buffer := range buffers {
			if !input.Quiet && len(input.Paths) > 1 && i == 0 && key > 0 {
				output.Write("\n==>%v<==\n", input.Paths[key])
				key++
			} else if key == 0 { key++ }

			if buffer.Done() {
				i = 0
				continue
			}

			output.Write("%s", buffer.Bytes())
			if i == 0 { i++ }
		}
	})

	return
}
