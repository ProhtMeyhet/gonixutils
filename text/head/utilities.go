package head

import(

	"github.com/ProhtMeyhet/libgosimpleton/iotool"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func prepareFileHelper(input *Input, output abstract.OutputInterface, exitCode *uint8) (helper *iotool.FileHelper) {
	helper = iotool.ReadOnly().ToggleFileAdviceReadSequential()
	helper.SetStdinStdout(abstract.STDIN_TOKEN, input.Stdin, input.Stdout)
	if input.NoCache { helper.ToggleFileAdviceDontNeed() }
	helper.SetE(func(name string, e error) {
		output.WriteEMessage(e, " on '%v'", name)
		if iotool.IsNotExist(e) {
			*exitCode = uint8(abstract.ERROR_FILE_NOT_FOUND)
		} else {
			*exitCode = uint8(abstract.ERROR_READING)
		}
	}); return
}
