package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/ls"
)

func main() {
	flags := NewLsFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Ls(input)))
}
