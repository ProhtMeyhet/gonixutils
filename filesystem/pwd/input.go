package pwd

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PathList	[]string

	// use PWD from environment, even if it contains symlinks
	// logical		bool

	// avoid all symlinks
	// physical	bool
}
