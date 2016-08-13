package sleep

import(
	"time"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Durations	[]time.Duration
	Until		time.Time
	// for verbose output
	Format		string

	DoUntil		bool
}
