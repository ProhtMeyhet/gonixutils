package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/user/whoami"
)

func main() {
	flags := NewWhoAmIFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.WhoAmI(input)))
}
