package main

import(
	"os"

	"github.com/ProhtMeyhet/gonixutils/text/echo"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type EchoFlags struct {
	abstract.Flags
	echo.Input

	unparsedType string
}

func NewEchoFlags() *EchoFlags {
	return &EchoFlags{}
}

func (flags *EchoFlags) GetInput() *echo.Input {
	return &flags.Input
}

func (flags *EchoFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("e", "escapes", "enable interpretation of backslash escapes", true)
	optarg.Add("n", "no-newline", "do not output trailing newline", false)
	optarg.Add("w", "to-stderr", "print to STDERR instead of STDOUT", false)

	// be compatible with gnu echo
	optarg.Add("E", "no-escapes", "disable interpretation of backslash escapes", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "e":
			flags.Escapes = option.Bool()
		case "E":
			flags.Escapes = !option.Bool()
		case "n":
			flags.NoNewLine = option.Bool()
		case "w":
			flags.PrintToStderr = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.Arguments = append(flags.Arguments, name)
	}
}

func (flags *EchoFlags) Validate() {
	flags.Flags.Validate()
}
