package main

import(
	"fmt"

	"os"
	"strings"
	"time"

	"github.com/ProhtMeyhet/gonixutils/miscellaneous/sleep"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type SleepFlags struct {
	abstract.Flags
	sleep.Input

	unparsedUntil string
}

func NewSleepFlags() *SleepFlags {
	return &SleepFlags{}
}

func (flags *SleepFlags) GetInput() *sleep.Input {
	return &flags.Input
}

func (flags *SleepFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("u", "until", "sleep until this clocktime (example: 1:00am)", "")

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "u":
			flags.DoUntil = option.String() != ""
			flags.unparsedUntil = option.String()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
	again:
		if duration, e := time.ParseDuration(name); e != nil {
			if strings.HasPrefix(e.Error(), "time: missing unit in duration") {
				name += "s"; goto again
			}
			fmt.Fprintf(os.Stderr, "couldn't parse duration '%v': '%v' -- ignoring\n", name, e)
		} else {
			flags.Durations = append(flags.Durations, duration)
		}
	}

	if flags.DoUntil {
		format, until, e := sleep.ParseUntil(flags.unparsedUntil)
		if e == nil {
			flags.Until = until
			flags.Format = format
		} else {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(sleep.EXIT_UNTIL_INVALID)
		}
	}
}

func (flags *SleepFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.Durations) == 0 && !flags.DoUntil {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
}

