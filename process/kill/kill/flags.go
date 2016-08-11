package main

import(
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
	"github.com/ProhtMeyhet/gonixutils/process/kill"

	"github.com/jteeuwen/go-pkg-optarg"
)

// cli flags
type flagConfig struct {
	abstract.Flags
	kill.Input
}

// new cli flags
func NewFlagConfig() *flagConfig {
	return &flagConfig{ }
}

// get our Input
func (flags *flagConfig) GetInput() *kill.Input {
	return &flags.Input
}

// well, duh
func (flags *flagConfig) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("f", "force", "one second after sending SIGQUIT send SIGKILL. please don't.", false)
// TODO	optarg.Add("i", "interactive", "one second after sending SIGQUIT send SIGSTOP and wait for user input", false)
	optarg.Add("n", "intervall", "intervall between the signals in seconds", 1)

	// -l, --list not implemented as this is for the man page
	// -L, --table not implemented; same reason
	// -p is deprecated

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
			case "f":
				flags.Force = option.Bool()
			case "i":
				flags.Interactive = option.Bool()
			case "n":
				flags.Intervall = option.Uint()
		}
	}

	flags.Help, flags.Verbose, flags.Version = flags.ParseFinally()

	for _, process := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end..
		if process == "" { continue }

		processId, e := strconv.Atoi(process)
		abstract.ExitOnError(e, os.Stderr, abstract.ERROR_PARSING, "'%v' arguments must be process or job IDs", process)
		if processId == 1 && !flags.Force { // remember 1 is the pid of init?
			fmt.Fprintf(os.Stderr, "good night sweet prince! you've got 5 seconds to hit CTRL + C!\n")
			time.Sleep(5 * time.Second)
		}
		flags.ProcessIds = append(flags.ProcessIds, processId)
	}
}

// validate. check if date is valid.
func (flags *flagConfig) Validate() {
	flags.Flags.Validate()

	if len(flags.ProcessIds) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
}
