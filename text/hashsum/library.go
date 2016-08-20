package hashsum

import(
	"encoding/hex"
	"hash"

	"github.com/ProhtMeyhet/libgosimpleton/simpleton"
	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)


// Those who do not understand UNIX are condemned to reinvent it, poorly.
//
//   — Henry Spencer

func doHash(input *Input, factory func() hash.Hash) (exitCode uint8) {
	output := abstract.NewOutput(input.Stdout, input.Stderr)
	helper := prepareFileHelper(input, output, &exitCode)
	exitCode = DoFromList(input, output, helper, factory, input.PathList...)
	output.Done(); output.Wait(); return
}

func Do(input *Input, output abstract.OutputInterface, helper *iotool.FileHelper, factory func() hash.Hash, paths <-chan string) (exitCode uint8) {
	parallel.OpenFilesDoWork(helper, paths, func(buffers chan *iotool.NamedBuffer) {
		hasher := factory()
		for buffered := range buffers {
			if buffered.Done() {
				output.Write("%v  %v\n", hex.EncodeToString(hasher.Sum(nil)), buffered.Name())
				hasher.Reset()
				continue
			}

			written, e := hasher.Write(buffered.Bytes())
			if written != buffered.Read() {
				output.WriteError("short write on hasher! output probably wrong!")
				exitCode = ERROR_HASH_FUNCTION
			}

			if output.WriteE(e) {
				exitCode = ERROR_HASH_FUNCTION
			}
		}
	})

	return
}

func DoFromList(input *Input, output abstract.OutputInterface, helper *iotool.FileHelper, factory func() hash.Hash, paths ...string) (exitCode uint8) {
	return Do(input, output, helper, factory, simpleton.StringListToChannel(paths...))
}
