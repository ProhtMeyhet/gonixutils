package hash

import(
	"encoding/hex"
	"hash"

//	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// use a hasher to hash a file or stdin
func doHash(input *Input, factory func() hash.Hash) (exitCode uint8) {
	work := parallel.NewStringFeeder(input.PathList)
	output := abstract.NewOutput(input.Stdout, input.Stderr)

	if input.VerboseLevel == 1 {
		output.WriteError("using %v workers for %v input(s).\n", work.Workers(), len(input.PathList))
	}

	work.Start(func() {
		exitCode = hash1(input, output, factory, work.Talk)
	})

	work.Wait()

	if input.VerboseLevel >= 2 {
		output.WriteError("used %v workers for %v input(s).\n", work.Workers(), len(input.PathList))
	}

	output.Done(); output.Wait()

	return
}

// hash one file, get it's path via list. when finished, reset and restart until list is closed.
func hash1(input *Input, output *abstract.Output, factory func() hash.Hash, list chan string) (exitCode uint8) {
	work := parallel.NewWork(1); hasher := factory()
	buffers := make(chan NamedBuffer, work.SuggestFileBufferSize())

	// open and read in one thread
	work.Feed(func() {
		defer close(buffers); helper := prepareFileHelper(input, output, &exitCode)
		ReadFiles(helper, buffers, list)
	})

	// hashing in another thread
	work.Run(func() {
		for buffered := range buffers {
			if buffered.done {
				output.Write("%v  %v\n", hex.EncodeToString(hasher.Sum(nil)), buffered.name)
				hasher.Reset()
				continue
			}

			written, e := hasher.Write(buffered.buffer[:buffered.read])
			if written != buffered.read {
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
