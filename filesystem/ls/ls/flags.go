package main

import(
//	"fmt"

	"os"

	"github.com/ProhtMeyhet/gonixutils/filesystem/ls"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type LsFlags struct {
	abstract.Flags
	ls.Input
}

func NewLsFlags() *LsFlags {
	return &LsFlags{}
}

func (flags *LsFlags) GetInput() *ls.Input {
	return &flags.Input
}

func (flags *LsFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("a", "all", "print everything, including . and ..", false)
	optarg.Add("c", "color", "set color true/false", true)
	optarg.Add("l", "list", "print files by line with details", false)
	optarg.Add("n", "new-line", "print each entry by line", false)
	optarg.Add("s", "sort", "set sort true/false, since multithreading the order without sort is likely random", true)
	optarg.Add("r", "reverse", "sort reverse", true)
	optarg.Add("R", "recursive", "recursivly list all directorys encountered", true)
	optarg.Add("u", "union", "treat all entries as if they are in the same directory", true)

//	optarg.Header("Behaviour changing:")
//	optarg.Add("r", "recursive", "create directorys recursive (works with files too)", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "a":
			flags.All = option.Bool()
		case "c":
			flags.NoColor = option.Bool()
		case "l":
			flags.Detail = option.Bool()
		case "n":
			flags.Lines = option.Bool()
		case "r":
			flags.SortReversed = option.Bool()
		case "R":
			flags.Recursive = option.Bool()
		case "s":
			flags.NoSort = option.Bool()
		case "u":
			flags.Union = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.Paths = append(flags.Paths, name)
	}

	if len(flags.Paths) == 0 {
		flags.Paths = append(flags.Paths, ".")
	}
}

func (flags *LsFlags) Validate() {
	flags.Flags.Validate()
}

