package head

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Paths		[]string

	Max		uint

	Quiet		bool
	NoCache		bool
	Bytes		bool
	Lines		bool
	Runes		bool
}
