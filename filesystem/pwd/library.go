package pwd

import(
	"os"
	"path"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/ProhtMeyhet/libgosimpleton/simpleton"
)

func Pwd(input *Input) (exitCode uint8) {
	output := abstract.NewOutput(input.Stdout, input.Stdout)
	exitCode = PrintWorkingDirectory(input, output)
	output.Done(); output.Wait(); return
}

func PrintWorkingDirectory(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	workingDirectory, e := os.Getwd(); if output.WriteE(e) { exitCode = abstract.FAILED; goto out }

	if len(input.PathList) > 0 {
		input.PathList = simpleton.SliceStringInsert(input.PathList, workingDirectory, 0)
		workingDirectory = path.Join(input.PathList...)
	}

	output.WriteLine(workingDirectory)

out:
	return
}
