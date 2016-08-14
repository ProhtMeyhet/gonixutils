package main

import(
	"fmt"

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
	optarg.Add("f", "file", "create file", false)
	optarg.Add("t", "temporary", "create a unique named temporary file under " +
				os.TempDir() + " and print its name. file input are prefixes.", false)
	optarg.Add("l", "link", "create link", false)

	optarg.Header("Behaviour changing:")
	optarg.Add("r", "recursive", "create directorys recursive (works with files too)", false)
	optarg.Add("s", "symbolic", "create soft/symbolic link. implies --link", false)
// TODO optarg.Add("p", "permissions", "set create permissions", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "f":
			flags.File = option.Bool()
		case "t":
			flags.Temporary = option.Bool()
		case "l":
			flags.Link = option.Bool()

		case "r":
			flags.Recursive = option.Bool()
		case "s":
			flags.Link = option.Bool()
			flags.Symbolic = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.PathList = append(flags.PathList, name)
	}
}

func (flags *LsFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.PathList) == 0 && !flags.Temporary {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}

	// TODO more checking
	if flags.Temporary && flags.Recursive {
		fmt.Fprintf(os.Stderr, "%v\n", "cannot recursivly create Temporary file!")
		os.Exit(abstract.ERROR_INVALID_ARGUMENT)
	}
}

