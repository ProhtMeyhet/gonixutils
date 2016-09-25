package ps

import(
	"strconv"

	"github.com/ProhtMeyhet/libgosimpleton/parallel"
	"github.com/ProhtMeyhet/libgosimpleton/system"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func Ps(input *Input) (exitCode uint8) {
	output := abstract.NewTabbedOutput(input.Stdout, input.Stderr)
	exitCode = Processes(input, output)
	output.Done(); output.Wait(); return
}

func Processes(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	if len(input.ProcessIds) == 0 && len(input.Processes) == 0 {
		if input.IsPosix {
			return OverviewPosix(input, output)
		} else {
			return Overview(input, output)
		}
	}

	ListProcessesById(input, output, input.ProcessIds...)
	ListProcessesByName(input, output, input.Processes...)
	return
}

func Overview(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	myProcesses := system.FindMyProcesses()
	for _, process := range myProcesses {
		out(input, output, "", process)
	}

	output.WriteFormatted("\nTotal: %v", len(myProcesses))

	return
}

func OverviewPosix(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	return
}

// FIXME remove repeated arguments (eg giving firefox twice or more)
func ListProcessesByName(input *Input, output abstract.OutputInterface, processes ...string) {
	if len(processes) == 0 { return }

	work := parallel.NewStringsFeeder(processes...)

	work.Start(func() {
		for name := range work.Talk {
			switch {
			case input.Oldest:
				process := system.FindOldestProcessByName(name)
				out(input, output, name, process)
			case input.Youngest:
				process := system.FindYoungestProcessByName(name)
				out(input, output, name, process)
			default:
				for _, process := range system.FindProcessesByName(name) {
					out(input, output, name, process)
				}
			}
		}
	})

	work.Wait()
}

// FIXME remove repeated arguments (eg giving pid 1 twice or more)
func ListProcessesById(input *Input, output abstract.OutputInterface, processIds ...uint64) {
	if len(processIds) == 0 { return }

	work := parallel.NewUints64Feeder(processIds...)

	work.Start(func() {
		for pid := range work.Talk {
			process, e := system.FindProcess(pid); if output.WriteE(e) { continue }
			out(input, output, strconv.FormatUint(pid, 10), process)
		}
	})

	work.Wait()
}

func out(input *Input, output abstract.OutputInterface, name string, process *system.ProcessInfo) {
	if process == nil {
		if input.Verbose {
			output.WriteFormatted("process '%v' not found.", name)
		}
		return
	}
	if input.Dump {
		output.WriteFormatted("%s\n", process.String())
	} else {
		output.WriteFormatted(decorate(input), process.Id(), process.Name())
	}
}

func decorate(input *Input) string {
	if input.NoColor {
		return "%v %v"
	}

	return abstract.TERMINAL_COLOR_GREEN + "%v" + abstract.TERMINAL_COLOR_RESET + " %v"
}
