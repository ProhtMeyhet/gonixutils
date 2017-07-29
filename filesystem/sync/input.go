package sync

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PathList	[]string

	// sync only file data, no unneeded metadata
	Data		bool

	// sync file
	File		bool

	// [not implemented] sync the file systems that contain the files
	FileSystem	bool
}
