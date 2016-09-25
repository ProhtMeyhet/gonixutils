package main

import(
	"fmt"

	"os"

	"github.com/ProhtMeyhet/gonixutils/filesystem/cp"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type CpFlags struct {
	abstract.Flags
	cp.Input
}

func NewCpFlags() *CpFlags {
	return &CpFlags{}
}

func (flags *CpFlags) GetInput() *cp.Input {
	return &flags.Input
}

func (flags *CpFlags) Parse() {
// cp [-prfvi] source... destination

	optarg.Header("Options for " + os.Args[0])
	optarg.Add("r", "recursive", "copy directories recursively", false)
	optarg.Add("i", "interactive", `prompt before overwrite`, false)
//	optarg.Add(false, `prompt before overwrite`, "i", func(option optarg.Option) {
//		flags.Interactive = option.Bool()
//	}, "interactive")

	optarg.Header("Behaviour changing:")
	optarg.Add("f", "force", `overwrite existing file`, false)
	optarg.Add("p", "preserve-execute", `also preserve the execute bit`, false)
	optarg.Add("o", "preserve-owner", `also preserve owner and owning group`, false)
	optarg.Add("t", "no-preserve", `do not try to preserve any file modes and use defaults`, false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "i":
			flags.Interactive = option.Bool()
		case "f":
			flags.Force = option.Bool()
		case "o":
			flags.PreserveOwner = option.Bool()
		case "p":
			flags.PreserveExecute = option.Bool()
		case "r":
			flags.Recursive = option.Bool()
		case "t":
			flags.NoPreserve = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.Paths = append(flags.Paths, name)
	}

	if len(flags.Paths) > 1 {
		flags.Destination = flags.Paths[len(flags.Paths) -1]
		flags.Paths = flags.Paths[:len(flags.Paths)-1]
	}

	if flags.VerboseLevel >= 5 {
		fmt.Printf("destination: %v\n", flags.Destination)
		fmt.Printf("files to copy: %v\n", flags.Paths)
	}
}

func (flags *CpFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.Paths) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
}
