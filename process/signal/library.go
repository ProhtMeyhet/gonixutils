package signal

import(
	"fmt"
	"os"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func Signal(input *Input) (exitCode uint8) {
	for _, processId := range input.ProcessIds {
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

	return
}
