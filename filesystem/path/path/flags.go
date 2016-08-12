package main

import(
	"fmt"

	"os"
	"strconv"

	"github.com/ProhtMeyhet/gonixutils/filesystem/path"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type PathFlags struct {
	abstract.Flags
	path.Input

	workingDirectoryE error
}

func NewPathFlags() *PathFlags {
	return &PathFlags{}
}

func (flags *PathFlags) GetInput() *path.Input {
	return &flags.Input
}

func (flags *PathFlags) Parse() {
	optarg.Header("Options for " + os.Args[0] + " ...DIR|FILE")
	optarg.Add("r", "relative", "print relative part from current working directory to DIR|FILE", false)
	optarg.Add("b", "basename", "strip directory", false)
	optarg.Add("d", "directory", "print directory name (dirname)", false)
	optarg.Add("e", "extension", "print the part beginning at the last dot. this is usually the file extension", false)
	optarg.Add("l", "list", "split by list separator (Unix ':')", false)
	optarg.Add("c", "clean", "return the shortest path name equivalent to path", false)
	optarg.Add("j", "join", "join all parts with directory seperator", false)
	optarg.Add("t", "split", "split path by directory separator and print each part by line", false)

	optarg.Header("Test options")
	optarg.Add("i", "is-absolute", "exit with " + strconv.Itoa(abstract.SUCCESS) + " if path is absolute, exit " +
				strconv.Itoa(abstract.FAILED) + " if not.", false)
	optarg.Add("m", "match", "exit with " + strconv.Itoa(abstract.SUCCESS) + " if pattern matched, exit " +
				strconv.Itoa(abstract.FAILED) + " if not.", "")

	optarg.Header("Behaviour options")
	optarg.Add("w", "working-directory", "overwrite working directory", "")
	optarg.Add("p", "prefix", "remove this prefix", "")
	optarg.Add("s", "suffix", "remove this suffix", "")

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
			case "r":
				flags.Relative = option.Bool()
			case "b":
				flags.Basename = option.Bool()
			case "d":
				flags.Directory = option.Bool()
			case "e":
				flags.Extension = option.Bool()
			case "l":
				flags.List = option.Bool()
			case "c":
				flags.Clean = option.Bool()
			case "j":
				flags.Join = option.Bool()
			case "m":
				flags.Pattern = option.String()
				flags.Match = true
			case "t":
				flags.Split = option.Bool()

			case "i":
				flags.IsAbsolute = option.Bool()

			case "w":
				flags.WorkingDirectory = option.String()
			case "p":
				flags.Prefix = option.String()
			case "s":
				flags.Suffix = option.String()
		}
	}

	if flags.WorkingDirectory == "" {
		flags.WorkingDirectory, flags.workingDirectoryE = os.Getwd()
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.PathList = append(flags.PathList, name)
	}
}

// TODO validate options that can't be used with other options
func (flags *PathFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.PathList) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}

	if flags.workingDirectoryE != nil {
		fmt.Println(flags.workingDirectoryE)
		os.Exit(abstract.ERROR_WORKING_DIRECTORY)
	}
}

