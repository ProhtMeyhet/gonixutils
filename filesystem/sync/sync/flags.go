package main

import(
//	"fmt"

	"os"

	"github.com/ProhtMeyhet/gonixutils/filesystem/sync"
	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/jteeuwen/go-pkg-optarg"
)

type SyncFlags struct {
	abstract.Flags
	sync.Input
}

func NewSyncFlags() *SyncFlags {
	return &SyncFlags{}
}

func (flags *SyncFlags) GetInput() *sync.Input {
	return &flags.Input
}

func (flags *SyncFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("d", "data", "sync data of [FILE] without metadata", false)
	optarg.Add("f", "file", "sync [FILE]", true)
	// TODO
	// optarg.Add("r", "recursive", "if [FILE] is a directory, recursivly sync each file", true)
	// optarg.Add("s", "filesystem", "sync filesystem of [FILE]", true)

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "d":
			flags.Data = option.Bool()
		case "f":
			flags.File = option.Bool()
		case "s":
			flags.FileSystem = option.Bool()
		}
	}

	flags.Help, flags.Version, flags.Verbose, flags.VerboseLevel = flags.ParseFinally()

	for _, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		flags.PathList = append(flags.PathList, name)
	}
}

func (flags *SyncFlags) Validate() {
	flags.Flags.Validate()
}

