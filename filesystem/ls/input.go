package ls

import(
	"os"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Paths		[]string

	// ls -l -- long. print files in detail per line
	Detail		bool

	// print all, including . and .. - implies Hidden
	All		bool

	// print hidden files
	Hidden		bool

	// print lines
	Lines		bool

	// dont sort
	NoSort		bool

	// sort reversed
	SortReversed	bool

	// turn off colors
	NoColor		bool

	// no decoration whatsoever
	NoDecoration	bool

	// breaks input a bit, but these functions are set corresponding to input values, so what.
	decorate	func(string, os.FileInfo) string
	writeEntry	func(input *Input, output abstract.OutputInterface, entry *Entry)
}
