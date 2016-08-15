package ls

import(


	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func Ls(input *Input) (exitCode uint8) {
	work := parallel.NewStringFeeder(input.Paths)
	output := abstract.NewSortedTabbedOutput(input.Stdout, input.Stderr)
	if input.Lines {
		output = abstract.NewOutput(input.Stdout, input.Stderr)
	} else if input.NoSort {
		output = abstract.NewTabbedOutput(input.Stdout, input.Stderr)
	}

	return list(input, output, work)
}
