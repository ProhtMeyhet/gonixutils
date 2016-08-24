package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/text/echo"
)

func main() {
	flags := NewEchoFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Echo(input)))
}
