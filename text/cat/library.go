package cat

import(


	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// unix cat
func Cat(input *Input) (exitCode uint8) {
	output := abstract.NewOutput(input.Stdout, input.Stderr)
	helper := prepareFileHelper(input, output, &exitCode)

	exitCode = Concat(input, output, helper)

	output.Done(); output.Wait()
	return
}

func Concat(input *Input, output abstract.OutputInterface, helper *iotool.FileHelper) (exitCode uint8) {
	key := 0
	parallel.ReadFilesSequential(helper, input.Paths, func(buffers chan *iotool.NamedBuffer) {
		i := 0
		if input.Verbose && len(input.Paths) > 1 && key == 0 { output.Write("==>%v<==\n", input.Paths[0]) }
		for buffer := range buffers {
			if input.Verbose && len(input.Paths) > 1 && i == 0 && key > 0 {
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
