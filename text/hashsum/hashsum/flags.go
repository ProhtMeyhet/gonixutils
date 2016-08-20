package main

import(
	"fmt"

	"os"
	"strings"

	"github.com/ProhtMeyhet/gonixutils/text/hashsum"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type HashFlags struct {
	abstract.Flags
	hashsum.Input

	unparsedType string
}

func NewHashFlags() *HashFlags {
	return &HashFlags{}
}

func (flags *HashFlags) GetInput() *hashsum.Input {
	return &flags.Input
}

func (flags *HashFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
//	optarg.Add("s", "salt", "first salt with this salt", "")
//	optarg.Add("r", "random-salt", "add a random salt of this length and print it after file name", 24)
	optarg.Add("c", "compare", "read sums from the FILEs and check them", false)
	optarg.Add("i", "idiot", "set the fail threshold for compare, when the hash function is detected as wrong", 10)
	optarg.Add("q", "quiet", "only print errors", false)
//	optarg.Add("w", "workers", "number of workers. if 0, the minimum of number of cpus * 2 or the number of input files", 0)
	if abstract.SET_FILE_ADVICE_DONTNEED {
		optarg.Add("x", "no-cache", "try not to leave caches behind (set file advice dont need)", false)
	}

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "x":
			flags.NoCache = option.Bool()
		case "c":
			flags.Compare = option.Bool()
		case "i":
			flags.Idiot = option.Uint()
		case "q":
			flags.Quiet = option.Bool()
		case "w":
			flags.NumberOfWorkers = option.Uint()
			if flags.NumberOfWorkers > abstract.GOROUTINE_LIMIT {
				flags.NumberOfWorkers = abstract.GOROUTINE_LIMIT
			}
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for key, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		if key == 0 {
			flags.unparsedType = strings.ToLower(name)
		} else {
			flags.PathList = append(flags.PathList, name)
		}
	}

	// read from stdin
	if len(flags.PathList) == 0 { flags.PathList = append(flags.PathList, "-") }

	flags.Type = hashsum.ParseType(flags.unparsedType)

	if flags.Idiot == 0 {
		flags.Idiot = 10
	}
}

func (flags *HashFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.PathList) > 0 && flags.Type == hashsum.NONE {
		fmt.Fprintf(os.Stderr, "%v is not a valid hashtype!\nexample usage: %v sha1 ...FILE\n", flags.unparsedType, os.Args[0])
		os.Exit(abstract.ERROR_INVALID_ARGUMENT)
	}

/*
	if len(flags.PathList) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}*/
}

