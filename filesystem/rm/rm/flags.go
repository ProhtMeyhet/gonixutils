package main

import(
	"os"

	"github.com/ProhtMeyhet/gonixutils/filesystem/rm"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type RmFlags struct {
	abstract.Flags
	rm.Input
}

func NewRmFlags() *RmFlags {
	return &RmFlags{}
}

func (flags *RmFlags) GetInput() *rm.Input {
	return &flags.Input
}

func (flags *RmFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("f", "force", "ignore if files do not exist, never prompt.", false)
	optarg.Add("R", "recursive", "remove directories and their contents recursively.", false)
	optarg.Add("r", "recursive", "remove directories and their contents recursively.", false)
	optarg.Add("i", "interactive", "prompt before each removal. implies --recursive.", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
			case "f":
				flags.Force = option.Bool()
			case "R":
				fallthrough
			case "r":
				flags.Recursive = option.Bool()
			case "i":
				// implies recursive!
				flags.Interactive = option.Bool()
				flags.Recursive = option.Bool()
		}
	}

	flags.Help, flags.Verbose, flags.Version = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.PathList = append(flags.PathList, name)
	}
}

func (flags *RmFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.PathList) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
}

