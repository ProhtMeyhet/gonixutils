package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/process/ps"
)

func main() {
	flags := NewPsFlags()
	flags.Parse()
	flags.Validate()

	flags.InitCli()

	exitCode := gonixutils.Ps(flags.GetInput())
	os.Exit(int(exitCode))
}
