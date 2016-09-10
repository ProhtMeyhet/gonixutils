package kill

import(
	"os"
	"syscall"
	"time"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
	"github.com/ProhtMeyhet/gonixutils/process/signal"
)

// in one second intervalls send first SIGTERM, then SIGINT, then SIGQUIT
// with input.Force, send SIGKILL after SIGQUIT
// with input.Interactive, send SIGSTOP and wait for user input
func Kill(input *Input) (exitCode uint8) {
	for _, processId := range input.ProcessIds {
		process, e := os.FindProcess(processId)
		if abstract.PrintErrorWithError(e, input.Stderr, "error: %v") {
			exitCode = signal.ERROR_FINDING_PROCESS
			continue
		}

		// request termination
		if Signal(input, process, syscall.SIGTERM) { continue }

		// request termination; terminal quit signal
		if Signal(input, process, syscall.SIGINT) { continue }

		// request termination and dump core
		if Signal(input, process, syscall.SIGQUIT) { continue }

		// kill cannot be caught nor ignored
		if input.Force {
			process.Signal(syscall.SIGKILL)
		} /* TODO else if input.Interactive {
			signal(input, process, syscall.SIGSTOP)
		}*/

		abstract.PrintError(`process '%v' did not respond to SIGTERM, SIGINT nor SIGQUIT. ` +
					`i suggest deleting the non conformant binary.`, processId)
	}

	return
}

// return true if to continue, false if not
func Signal(input *Input, process *os.Process, signal os.Signal) bool {
	if e := process.Signal(signal); e != nil {
		if e.Error() != "os: process already finished" {
			abstract.PrintErrorWithError(e, input.Stderr, "")
		}
		return true // error happend, cannot do anything else
	}

	return poll(process)
}

// in 100 millisecond intervalls test if process is still there
func poll(process *os.Process) bool {
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		// FIXME parse error message
		_, e := os.FindProcess(process.Pid); if e != nil {
			return true
		}
	}

	return false
}
