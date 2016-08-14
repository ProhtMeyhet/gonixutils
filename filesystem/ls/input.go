package mk

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Paths		[]string

	Sort		bool
	SortReversed	bool
}
