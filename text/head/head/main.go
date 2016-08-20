package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/text/head"
)

func main() {
	flags := NewHeadFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Head(input)))
}
