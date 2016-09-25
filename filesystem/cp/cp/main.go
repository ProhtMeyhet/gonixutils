package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/cp"
)

func main() {
	flags := NewCpFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Cp(input)))
}
