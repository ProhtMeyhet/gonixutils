package hashsum

import(
	"bufio"
	"bytes"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func Compare(input *Input) (exitCode uint8) {
	work := parallel.NewStringFeeder(input.PathList)
	output := abstract.NewOutput(input.Stdout, input.Stderr)

// TODO
//	if input.OrderedOutput {
//		exitCode = compareOrdered(work, output)
//	} else {
		exitCode = compare(input, output, work)
//	}

	output.Done(); output.Wait()

	return
}

func compare(input *Input, output abstract.OutputInterface, work *parallel.WorkString) (exitCode uint8) {
	work.Start(func() {
		exitCode = compare1(input, output, work.Talk)
	})

	work.Wait()

	return
}


func compare1(input *Input, output abstract.OutputInterface, list chan string) (exitCode uint8) {
	notFound, unequalHashes, wrongHashFunction, hashlen, expectedHashLen := 0, 0, uint(0), 0, HashLen(input.Type)
	reader, helper := &bufio.Reader{}, prepareFileHelper(input, output, &exitCode)
	work := parallel.NewWork(0); talk := make(chan HashPath, work.SuggestBufferSize(0))

	work.Feed(func() {
		for path := range list {
			handler, e := iotool.Open(helper, path); if e != nil { continue }
			reader.Reset(handler); scanner := bufio.NewScanner(reader)
			scanner.Split(bufio.ScanWords); hashPath := HashPath{}

			for i := 1; scanner.Scan(); i++ {
				if i % 2 != 0 {
					hashPath.Hash = scanner.Bytes()
				} else {
					hashPath.Path = scanner.Bytes()
					if uint(len(hashPath.Hash)) != expectedHashLen {
						// no lock required (yet)
						wrongHashFunction++
						output.WriteError("%s, wrong hash function!\n", hashPath.Path)

						// give up
						if wrongHashFunction == input.Idiot {
							output.WriteError("file %s is probably not valid. giving up!\n", path)
							exitCode = ERROR_INVALID_FILE_FORMAT
							break
						}
					} else {
						talk <-hashPath
					}
				}
			}
		}; close(talk)
	})

	work.Start(func() {
		hash1List := make(chan string, work.SuggestBufferSize(0))
		in := bytes.NewBuffer(make([]byte, 0)); buffered := abstract.NewOutput(in, input.Stderr)
		iwork := parallel.NewWork(work.Workers()); hashPaths := make(map[string]HashPath)

		iwork.Feed(func() {
			for hashPath := range talk {
				hash1List <-string(hashPath.Path)
				hashPaths[string(hashPath.Path)] = hashPath
				hashlen = len(hashPath.Hash)
			}; close(hash1List)
		})

		iwork.Start(func() {
			Do(input, buffered, Factory(input), hash1List)
		})

		// since the output is unordered, wait until everything is finished before comparing
		// comparing is very quick, even with many files, as it's still only a few bytes
		iwork.Wait(); buffered.Done(); buffered.Wait()

		hash, path := make([]byte, 0), ""
		scanner := bufio.NewScanner(in); scanner.Split(bufio.ScanWords)
		for i := 1; scanner.Scan(); i++ {
			if i % 2 != 0 {
				hash = scanner.Bytes()
			} else {
				path = scanner.Text()

				// the if only guards nil panics
				if value, ok := hashPaths[path]; ok {
					if !bytes.Equal(hash, value.Hash) {
						output.WriteError("%s: FAILED!\n", path)
						work.Lock(); unequalHashes++; work.Unlock()
					} else if !input.Quiet {
						output.Write("%s: OK\n", path)
					}
					delete(hashPaths, path)
				}
			}
		}

		if len(hashPaths) > 0 {
			exitCode = abstract.FAILED
			work.Lock(); notFound += len(hashPaths); work.Unlock()
			for path, _ := range hashPaths {
				output.WriteError("%s: FAILED! file not found.\n", path)
			}
		}
	})

	work.Wait()

	if wrongHashFunction > 0 {
		output.WriteError(`hashsum: WARNING: on %v listed file(s) the hash size was unexpected! ` +
					`probably wrong hash function!` + "\n", wrongHashFunction)
	}

	if notFound > 0 {
		output.WriteError("hashsum: WARNING: %v listed file(s) could not be read or not found!\n", notFound)
	}

	if unequalHashes > 0 {
		output.WriteError("hashsum: WARNING: %v computed checksum(s) did NOT match!\n", unequalHashes)
	}

	return
}

type HashPath struct {
	Hash, Path []byte
}
