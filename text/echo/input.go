package echo

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Arguments	[]string

	NoNewLine	bool
	Escapes		bool
	PrintToStderr	bool
//	NumberLines	bool
}
