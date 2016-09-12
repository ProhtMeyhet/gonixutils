package cat

import(
	"fmt"
	"io"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// unix cat
func Cat(input *Input) (exitCode uint8) {
	output := abstract.NewOutput(input.Stdout, input.Stderr)
	if input.Verbose { output.TogglePrintSubBufferNames() }
	helper := prepareFileHelper(input, output, &exitCode)

	e := WriteFilesToOutput(output, helper, input.Paths...); if e != nil && exitCode == 0 {
		// TODO parse the error and set exitCode accordingly
		exitCode = abstract.ERROR_UNHANDLED
	}

	output.Done(); output.Wait(); return
}

func WriteFilesToOutput(mainOutput abstract.OutputInterface, helper *iotool.FileHelper, paths ...string) (e error) {
	return WriteFilesFilteredToOutput(mainOutput, helper, nil, paths...)
}

func WriteFilesFilteredToOutput(mainOutput abstract.OutputInterface, helper *iotool.FileHelper,
		filter func(io.Reader) io.Reader, paths ...string) (e error) {
	output := mainOutput
	parallel.ReadFilesSequential(helper, paths, func(buffered *iotool.NamedBuffer) {
		var filtered io.Reader
		// use a subbuffer if required
		if output.PrintSubBufferNames() {
			output = mainOutput.NewSubBuffer(fmt.Sprintf("==>%v<==\n", buffered.Name()), 0)
		}

		// apply filter
		if filter != nil { filtered = filter(buffered) }
		// if the filter func failed, use the buffered
		if filtered == nil { filtered = buffered }

		// copy and close
		_, e = io.Copy(output, filtered)
		if output.PrintSubBufferNames() { output.Done() }
	}).Wait(); return
}
