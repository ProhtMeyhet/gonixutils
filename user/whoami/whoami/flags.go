package main

import(
	"os"

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
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("a", "all", "print all available information formatted", false)
	optarg.Add("i", "id", "print my user id", false)
	optarg.Add("g", "gid", "print my main group id", false)
	optarg.Add("n", "name", "print my real or display name", false)
	optarg.Add("u", "username", "print my username (default if no other option)", false)
	optarg.Add("h", "home", "print my home directory (if i have one)", false)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "a":
			flags.PrintAll = option.Bool()
		case "i":
			flags.PrintUid = option.Bool()
		case "g":
			flags.PrintGid = option.Bool()
		case "n":
			flags.PrintName = option.Bool()
		case "u":
			flags.PrintUsername = option.Bool()
		case "h":
			flags.PrintHome = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()
}

func (flags *WhoAmIFlags) Validate() {
	flags.Flags.Validate()
}
