package hash

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PathList	[]string
	Type		Type

	Salt		[]byte

	Quiet		bool
	NoCache		bool
	NumberOfWorkers uint

	currentPath	string
}
