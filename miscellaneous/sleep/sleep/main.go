package main

import(
	"os"

	gonixutils "github.com/ProhtMeyhet/gonixutils/miscellaneous/sleep"
)

func main() {
	flags := NewSleepFlags(); flags.Parse(); flags.Validate()

	input := flags.GetInput(); input.InitCli()

	go gonixutils.SignalHandler(input)
	os.Exit(int(gonixutils.Sleep(input)))
}
