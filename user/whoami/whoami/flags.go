package main

import(
//	"os"

	"github.com/ProhtMeyhet/gonixutils/user/whoami"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type WhoAmIFlags struct {
	abstract.Flags
	whoami.Input

	unparsedType string
}

func NewWhoAmIFlags() *WhoAmIFlags {
	return &WhoAmIFlags{}
}

func (flags *WhoAmIFlags) GetInput() *whoami.Input {
	return &flags.Input
}

func (flags *WhoAmIFlags) Parse() {
//	optarg.Header("Options for " + os.Args[0])
//	optarg.Add("e", "escapes", "enable interpretation of backslash escapes", true)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
/*		switch option.ShortName {
		case "e":
			flags.Escapes = option.Bool()
		}*/
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()
}

func (flags *WhoAmIFlags) Validate() {
	flags.Flags.Validate()
}
