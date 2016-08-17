package ls

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// list filesystem entries
func Ls(input *Input) (exitCode uint8) {
	output := abstract.NewSortedTabbedOutput(input.Stdout, input.Stderr)
	if input.Lines {
		output = abstract.NewOutput(input.Stdout, input.Stderr)
	} else if input.NoSort && !input.Detail {
		output = abstract.NewTabbedOutput(input.Stdout, input.Stderr)
	}

	exitCode = List(input, input.Paths, output)

	output.Done(); output.Wait()

	return
}