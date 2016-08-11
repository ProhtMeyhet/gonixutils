package kill

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Force bool
	Intervall uint
	Interactive bool

	ProcessIds []int
}
