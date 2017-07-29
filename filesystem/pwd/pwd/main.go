package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/pwd"
)

func main() {
	flags := NewPwdFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Pwd(input)))
}
