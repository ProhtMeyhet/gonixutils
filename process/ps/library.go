package ps

import(
	"strconv"

	"github.com/ProhtMeyhet/libgosimpleton/parallel"
	"github.com/ProhtMeyhet/libgosimpleton/system/processes"

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

	output.ToggleLinesManual()
	ListProcessesById(input, output, input.ProcessIds...)
	ListProcessesByName(input, output, input.Processes...)
	return
}

func Overview(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	if input.Dump { output.ToggleLinesManual() }
	myProcesses := processes.FindMyAll()
	for process := range myProcesses {
		if input.Dump {
			output.WriteFormatted(decorate(input) + "\n", process.Id(), process.CommandLine())
		} else {
			output.WriteFormatted(decorate(input), process.Id(), process.Name())
		}
	}

	output.WriteFormatted("\nTotal: %v", len(myProcesses))

	return
}

func OverviewPosix(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	return
}

// FIXME remove repeated arguments (eg giving firefox twice or more)
func ListProcessesByName(input *Input, output abstract.OutputInterface, process ...string) {
	if len(process) == 0 { return }

	work := parallel.NewStringsFeeder(process...)

	work.Start(func() {
		for name := range work.Talk {
			switch {
			case input.Oldest:
				process := processes.FindOldestByName(name)
				out(input, output, name, process)
			case input.Youngest:
				process := processes.FindYoungestByName(name)
				out(input, output, name, process)
			default:
				for process := range processes.FindByName(name) {
					out(input, output, name, process)
				}
			}
		}
	})

	work.Wait()
}

// FIXME remove repeated arguments (eg giving pid 1 twice or more)
func ListProcessesById(input *Input, output abstract.OutputInterface, processId ...uint) {
	if len(processId) == 0 { return }

	work := parallel.NewUintsFeeder(processId...)

	work.Start(func() {
		for pid := range work.Talk {
			process, e := processes.Find(pid); if output.WriteE(e) { continue }
			out(input, output, strconv.FormatUint(uint64(pid), 10), process)
		}
	})

	work.Wait()
}

func out(input *Input, output abstract.OutputInterface, name string, process *processes.ProcessInfo) {
	if process == nil {
		if input.Verbose {
			output.WriteFormatted("process '%v' not found.", name)
		}
		return
	}
	if input.Dump {
		output.WriteFormatted("%s\n", process.String())
	} else {
		output.WriteFormatted(decorate(input) + "\n", process.Id(), process.Name())
	}
}

func decorate(input *Input) string {
	if input.NoColor {
		return "%v %v"
	}

	return abstract.TERMINAL_COLOR_GREEN + "%v" + abstract.TERMINAL_COLOR_RESET + " %v"
}
