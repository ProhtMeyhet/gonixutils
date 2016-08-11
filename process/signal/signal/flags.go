package main

import(
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
	"github.com/ProhtMeyhet/gonixutils/process/signal"

	"github.com/jteeuwen/go-pkg-optarg"
)

type flagConfig struct {
	abstract.Flags
	signal.Input

	unparsedSignal string
}

func NewFlagConfig() *flagConfig {
	return &flagConfig{ unparsedSignal: signal.DEFAULT_SIGNAL }
}

func (flags *flagConfig) GetInput() *signal.Input {
	return &flags.Input
}

func (flags *flagConfig) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("s", "signal", "The signal to send.  It may be given as a name or a number.", signal.DEFAULT_SIGNAL)
	optarg.Add("f", "force", "Ignore safety and do whatever you please.", false)

	// -l, --list not implemented as this is for the man page
	// -L, --table not implemented; same reason
	// -p is deprecated
	// TODO -a
	// TODO -q

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
			case "s":
				flags.unparsedSignal = option.String()
			case "f":
				flags.Force = option.Bool()
		}
	}

	flags.Help, flags.Verbose, flags.Version = flags.ParseFinally()

	if flags.unparsedSignal != signal.DEFAULT_SIGNAL {
		// first try to parse as a number
		parsedSignal, e := strconv.Atoi(flags.unparsedSignal)

		if e != nil {
			flags.Signal, e = signal.StringToSignal(flags.unparsedSignal)
			abstract.ExitOnError(e, os.Stderr, abstract.ERROR_PARSING,
						"'%v': invalid signal specification", flags.unparsedSignal)
		} else {
			flags.Signal, e = signal.IntToSignal(parsedSignal)
			abstract.ExitOnError(e, os.Stderr, abstract.ERROR_PARSING,
						"'%v': invalid signal specification", flags.unparsedSignal)
		}
	} else {
		flags.Signal = syscall.SIGUSR1
	}

	for _, process := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end..
		if process == "" { continue }

		processId, e := strconv.Atoi(process)
		abstract.ExitOnError(e, os.Stderr, abstract.ERROR_PARSING,
					"'%v' arguments must be process or job IDs", process)
		if processId == 1 && signal.IsTerminal(flags.Signal) && !flags.Force { // remember 1 is the pid of init?
			fmt.Fprintf(os.Stderr, "good night sweet prince! you've got 5 seconds to hit CTRL + C!\n")
			time.Sleep(5 * time.Second)
		}
		flags.ProcessIds = append(flags.ProcessIds, processId)
	}
}

func (flags *flagConfig) Validate() {
	flags.Flags.Validate()

	if len(flags.ProcessIds) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
}
