package main

import(
//	"fmt"

	"os"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
	"github.com/ProhtMeyhet/gonixutils/filesystem/pwd"

	"github.com/jteeuwen/go-pkg-optarg"
)

type PwdFlags struct {
	abstract.Flags
	pwd.Input
}

func NewPwdFlags() *PwdFlags {
	return &PwdFlags{}
}

func (flags *PwdFlags) GetInput() *pwd.Input {
	return &flags.Input
}

func (flags *PwdFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	// TODO
	// optarg.Add("l", "logical", "use PWD from environment, even if it contains symlinks", false)
	// optarg.Add("p", "physical", "avoid all symlinks", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
	/*	switch option.ShortName {
		case "l":
			flags.Logical = option.Bool()
		case "p":
			flags.Physical = option.Bool()
		}*/
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.PathList = append(flags.PathList, name)
	}
}

func (flags *PwdFlags) Validate() {
	flags.Flags.Validate()
}

