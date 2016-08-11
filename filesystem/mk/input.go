package mk

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PathList	[]string

	File		bool
	Link		bool
	Temporary	bool
	Recursive	bool
	Interactive	bool
	Symbolic	bool
}
