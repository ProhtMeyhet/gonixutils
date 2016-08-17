package hashsum

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PathList		[]string
	Type			Type

	Salt			[]byte

	Compare			bool
	OrderCompareOutput	bool
	Quiet			bool
	NoCache			bool
	NumberOfWorkers		uint
	Idiot			uint
}
