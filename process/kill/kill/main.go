package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/process/kill"
)

func main() {
	flags := NewFlagConfig()
	flags.Parse()
	flags.Validate()

	flags.InitCli()

	exitCode := gonixutils.Kill(flags.GetInput())
	os.Exit(int(exitCode))
}
