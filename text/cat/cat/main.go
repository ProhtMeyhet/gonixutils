package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/text/cat"
)

func main() {
	flags := NewCatFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Cat(input)))
}
