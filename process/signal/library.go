package signal

import(
	"fmt"
	"os"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/ProhtMeyhet/libgosimpleton/parallel"
)

func Signal(input *Input) (exitCode uint8) {
	work := parallel.NewIntsFeeder(input.ProcessIds...)

	work.Start(func() {
		for processId := range work.Talk {
			process, e := os.FindProcess(processId)
			if abstract.PrintErrorWithError(e, input.Stderr, "error: %v") {
				exitCode = abstract.ERROR_GENERIC
				continue
			}

			if input.Verbose { fmt.Fprintf(input.Stderr, "signaling: '%v' with %v\n",
							processId, SignalToString(input.Signal) ) }

			e = process.Signal(input.Signal)
			if abstract.PrintErrorWithError(e, input.Stderr, "error: %v") {
				exitCode = abstract.ERROR_GENERIC
			}
		}
	})

	work.Wait(); return
}
