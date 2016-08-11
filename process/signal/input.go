package signal

import(
	"syscall"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Force bool

	ProcessIds []int
	Signal syscall.Signal
}
