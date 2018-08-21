package main

import(
	"os"
	"strconv"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
	"github.com/ProhtMeyhet/gonixutils/process/ps"

	"github.com/jteeuwen/go-pkg-optarg"
)

// cli flags
type PsFlags struct {
	abstract.Flags
	ps.Input
}

// new cli flags
func NewPsFlags() *PsFlags {
	return &PsFlags{ }
}

// get our Input
func (flags *PsFlags) GetInput() *ps.Input {
	return &flags.Input
}

// well, duh
func (flags *PsFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("a", "all", "don't ignore other users processes", false)
	optarg.Add("e", "exact", "match only exactly this name", false)
	optarg.Add("d", "dump", "dump all information about process available", false)
	optarg.Add("n", "newest", "select only the newest (most recently started) of the matching processes", false)
	optarg.Add("o", "oldest", "select only the oldest (least recently started) of the matching processes", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
			case "l":
				flags.List = option.Bool()
			case "a":
				flags.All = option.Bool()
			case "d":
				flags.Dump = option.Bool()
			case "n":
				flags.Youngest = option.Bool()
			case "o":
				flags.Oldest = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, process := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end..
		if process == "" { continue }

		processId, e := strconv.ParseUint(process, 10, strconv.IntSize); if e == nil {
			flags.ProcessIds = append(flags.ProcessIds, uint(processId))
		} else {
			flags.Processes = append(flags.Processes, process)
		}
	}
}

// validate. check if date is valid.
func (flags *PsFlags) Validate() {
	flags.Flags.Validate()
/*
	if len(flags.ProcessIds) == 0 && len(flags.Processes) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
*/
}
