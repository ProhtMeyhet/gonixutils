package path

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	// must set!
	WorkingDirectory string

	// pattern for match
	Pattern		string

	PathList	[]string

	Clean		bool
	Extension	bool
	Basename	bool
	Directory	bool
	List		bool
	Join		bool
	Split		bool
	Match		bool
	Relative	bool

	IsAbsolute	bool

	Prefix		string
	Suffix		string
}
