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
		for buffer := range buffers {
			if len(input.Paths) > 1 && !input.Quiet && buffer.Next() {
				if key == 0 {
					output.Write("==>%v<==\n", input.Paths[0])
				} else {
					output.Write("\n==>%v<==\n", input.Paths[key])
				}
				key++
			}

			output.Write("%s", buffer.Bytes())
		}
	}).Wait()

	return
}
