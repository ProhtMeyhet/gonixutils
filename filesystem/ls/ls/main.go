package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/ls"
)

func main() {
	flags := NewMkFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Ls(input)))
}
