package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/text/hashsum"
)

func main() {
	flags := NewHashFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Hashsum(input)))
}
