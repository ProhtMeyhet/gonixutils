package ps

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Processes	[]string
	// FIXME strconv doesn't have support for uint only uint64
	ProcessIds	[]uint

	All		bool
	List		bool
	Exact		bool

	Oldest		bool
	Youngest	bool
	Dump		bool

	NoColor		bool

	// be posix, not gonix
	IsPosix		bool
}

func Posix() (input *Input) {
	input = &Input{}; input.Posix(); return
}

func (input *Input) Posix() {
	input.IsPosix = true
}
