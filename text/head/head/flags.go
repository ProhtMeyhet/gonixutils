package main

import(
	"os"

	"github.com/ProhtMeyhet/gonixutils/text/head"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type HeadFlags struct {
	abstract.Flags
	head.Input

	unparsedType string
}

func NewHeadFlags() *HeadFlags {
	return &HeadFlags{}
}

func (flags *HeadFlags) GetInput() *head.Input {
	return &flags.Input
}

func (flags *HeadFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("b", "bytes", "print the first NUM bytes", uint(100))
	// be compatible with unix head
	optarg.Add("c", "bytes", "print the first NUM bytes", uint(100))
	optarg.Add("n", "lines", "print the first NUM lines", uint(10))
	optarg.Add("r", "runes", "print the first NUM runes. A rune is one utf8 character", uint(100))
	optarg.Add("q", "quiet", "dont print headers giving file names", false)


	if abstract.SET_FILE_ADVICE_DONTNEED {
		optarg.Add("x", "no-cache", "try not to leave caches behind (set file advice dont need)", false)
	}

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "n":
			flags.Max = option.Uint()
		case "b", "c":
			flags.Bytes = true
			flags.Max = option.Uint()
		case "q":
			flags.Quiet = option.Bool()
		case "r":
			flags.Runes = true
			flags.Max = option.Uint()
		case "x":
			if abstract.SET_FILE_ADVICE_DONTNEED { flags.NoCache = option.Bool() }
		}
	}

	if flags.Max == 0 { flags.Max = 10 }

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

func (flags *HeadFlags) Validate() {
	flags.Flags.Validate()

/*
	if len(flags.Paths) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}*/
}

