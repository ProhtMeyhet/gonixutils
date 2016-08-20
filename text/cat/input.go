package cat

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Paths		[]string

	NoCache		bool
	NumberLines	bool
}
