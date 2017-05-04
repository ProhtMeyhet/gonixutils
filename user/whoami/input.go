package whoami

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	PrintAll	bool
	PrintUsername	bool
	PrintUid	bool
	PrintGid	bool
	PrintName	bool
	PrintHome	bool
}

func (input *Input) Verify() (e error) {
	// if no option given, print Username
	if !input.PrintUsername && !input.PrintAll &&
		!input.PrintUid && !input.PrintGid && !input.PrintName && !input.PrintHome {
		input.PrintUsername = true
	}

	return
}
