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
	hash.Input

	unparsedType string
}

func NewHashFlags() *HashFlags {
	return &HashFlags{}
}

func (flags *HashFlags) GetInput() *hash.Input {
	return &flags.Input
}

func (flags *HashFlags) Parse() {
	optarg.Header("Options for " + os.Args[0])
	optarg.Add("s", "salt", "first salt with this salt", "")
//	optarg.Add("r", "random-salt", "add a random salt of this length and print it after file name", 24)
	optarg.Add("q", "quiet", "only print errors", false)
	optarg.Add("w", "workers", "number of workers. if 0, the minimum of number of cpus * 2 or the number of input files", 0)
	if abstract.SET_FILE_ADVICE_DONTNEED {
		optarg.Add("x", "no-cache", "try not to leave caches behind (set file advice dont need)", false)
	}

	flags.AddGeneralOptions()

	for option := range optarg.Parse() {
		if flags.ParseOption(option) { continue }
		switch option.ShortName {
		case "x":
			flags.NoCache = option.Bool()
		case "q":
			flags.Quiet = option.Bool()
		case "w":
			flags.NumberOfWorkers = option.Uint()
			if flags.NumberOfWorkers > abstract.GOROUTINE_LIMIT {
				flags.NumberOfWorkers = abstract.GOROUTINE_LIMIT
			}
		}
	}

	flags.Help, flags.Verbose, flags.Version = flags.ParseFinally()

	for key, name := range optarg.Remainder {
		// bash adds sometimes an empty argument at the end...
		if name == "" { continue }
		if key == 0 {
			flags.unparsedType = strings.ToLower(name)
		} else {
			flags.PathList = append(flags.PathList, name)
		}
	}

	switch flags.unparsedType {
	case "sha512", "sha512sum":
		flags.Type = hash.SHA512
	case "md5", "md5sum":
		flags.Type = hash.MD5
	case "sha1", "sha1sum":
		flags.Type = hash.SHA1
	case "sha256", "sha256sum":
		flags.Type = hash.SHA256
	case "sha224", "sha224sum":
		flags.Type = hash.SHA224
	case "sha384", "sha384sum":
		flags.Type = hash.SHA384
	case "sha512a", "sha512asum", "sha512_224", "sha512224":
		flags.Type = hash.SHA512_224
	case "sha512b", "sha512bsum", "sha512_256", "sha512256":
		flags.Type = hash.SHA512_256
	case "adler32", "adler32sum":
		flags.Type = hash.ADLER32
	case "crc32", "crc32sum":
		flags.Type = hash.CRC32
	case "crc64", "crc64sum":
		flags.Type = hash.CRC64
	}
}

func (flags *HashFlags) Validate() {
	flags.Flags.Validate()

	if len(flags.PathList) > 0 && flags.Type == hash.NONE {
		fmt.Fprintf(os.Stderr, "%v is not a valid hashtype!\nexample usage: %v sha1 ...FILE\n", flags.unparsedType, os.Args[0])
		os.Exit(abstract.ERROR_INVALID_ARGUMENT)
	}

	if len(flags.PathList) == 0 {
		flags.Usage()
		os.Exit(abstract.ERROR_NO_INPUT)
	}
}

