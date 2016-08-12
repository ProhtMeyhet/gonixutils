package abstract

import(
	"os"

	"github.com/jteeuwen/go-pkg-optarg"

	"github.com/ProhtMeyhet/gonixutils/library/version"
)

type Flags struct {
	questionMarkCount uint
	HasVerbose, HasHelp, HasVersion bool
	// tells you the level of verbosiness the user wants
	FlagVerboseLevel uint
	Debug bool
}

func (flags *Flags) AddGeneralOptions() {
	optarg.Header("General Options")
	optarg.Add("?", "help", "Display this help and exit", false)
	optarg.Add("??", "version", "Output version information and exit", false)
	optarg.Add("???", "verbose", "Be verbose", false)

/*
	optarg.Add("#", "goroutines", "set runtime.GOMAXPROCS", 8)
	optarg.Add("#", "count", "count xxx", 10)
	optarg.Add("*", "log-to", "log target: file, stderr, syslog". "syslog")
	optarg.Add("@", "user", "only this user", "")
	optarg.Add("$", "dry-run", "do not change anything (no cost)", false)
	optarg.Add("$", "scripting", "Be verbose", false)
	optarg.Add("^", "not", "negate all options", false)
	optarg.Add("%", "human", "human readable; MB or %", false)
	optarg.Add("+", "print-all", "show every information there is", false)
*/

	if DEBUG {
		optarg.Add("!", "debug", "Debug this", false)
	}
}

func (flags *Flags) ParseOption(option *optarg.Option) bool {
	switch option.ShortName {
		case "?":
			flags.questionMarkCount++
			return true
		case "??":
			flags.HasVersion = true
			return true
		case "!":
			flags.Debug = true
			return true
	}

	return false
}

func (flags *Flags) ParseFinally() (bool, bool, bool, uint) {
	switch flags.questionMarkCount {
		case 1:
			flags.HasHelp = true
		case 2:
			flags.HasVersion = true
	}

	if flags.questionMarkCount >= 3 {
		flags.HasVerbose = true
		flags.FlagVerboseLevel = flags.questionMarkCount - 2
	}


	return flags.HasHelp, flags.HasVersion, flags.HasVerbose, flags.FlagVerboseLevel
}

func (flags *Flags) Parse() <-chan *optarg.Option {
	flags.AddGeneralOptions()
	return optarg.Parse()
}

func (flags *Flags) Validate() {
	if flags.HasVersion {
		version.Print()
		os.Exit(0)
	}

	if flags.HasHelp {
		flags.Usage()
		os.Exit(0)
	}
}

func (flags *Flags) Usage() {
	optarg.Usage()
}
