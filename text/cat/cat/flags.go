package main

import(
	"os"

	"github.com/ProhtMeyhet/gonixutils/text/cat"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type CatFlags struct {
	abstract.Flags
	cat.Input

	unparsedType string
}

func NewCatFlags() *CatFlags {
	return &CatFlags{}
}

func (flags *CatFlags) GetInput() *cat.Input {
	return &flags.Input
}

func (flags *CatFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
//	optarg.Add("n", "number", "number each output lines", false)

	if abstract.SET_FILE_ADVICE_DONTNEED {
		optarg.Add("x", "no-cache", "try not to leave caches behind (set file advice dont need)", false)
	}

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "n":
			flags.NumberLines = option.Bool()
		case "x":
			if abstract.SET_FILE_ADVICE_DONTNEED { flags.NoCache = option.Bool() }
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.Paths = append(flags.Paths, name)
	}

	if len(flags.Paths) == 0 {
		flags.Paths = append(flags.Paths, abstract.STDIN_TOKEN)
	}
}

func (flags *CatFlags) Validate() {
	flags.Flags.Validate()

/*
	if len(flags.Paths) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}*/
}
