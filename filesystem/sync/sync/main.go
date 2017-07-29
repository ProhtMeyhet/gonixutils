package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/filesystem/sync"
)

func main() {
	flags := NewSyncFlags()
	flags.Parse()
	flags.Validate()

	input := flags.GetInput()
	input.InitCli()

	os.Exit(int(gonixutils.Sync(input)))
}
