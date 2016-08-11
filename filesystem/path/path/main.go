package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/path"
)

func main() {
	flags := NewPathFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Path(input)))
}
