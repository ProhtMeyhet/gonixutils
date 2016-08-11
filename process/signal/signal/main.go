package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/process/signal"
)

func main() {
	flags := NewFlagConfig()
	flags.Parse()
	flags.Validate()

	flags.InitCli()

	exitCode := gonixutils.Signal(flags.GetInput())
	os.Exit(int(exitCode))
}
