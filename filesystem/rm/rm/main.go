package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/rm"
)

func main() {
	flags := NewRmFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Rm(input)))
}
