package rm

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PathList []string

	Force bool
	Recursive bool
	Interactive bool
}
