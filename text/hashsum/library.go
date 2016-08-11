package hash

import(
	"encoding/hex"
	"hash"

	"crypto/sha1"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
//	"hash/fnv"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// use a hasher to hash a file or stdin
func doHash(input *Input, factory func() hash.Hash) (exitCode uint8) {
	numberOfWorkers := input.NumberOfWorkers
	if numberOfWorkers == 0 {
		numberOfWorkers = parallel.SuggestNumberOfWorkers(uint(len(input.PathList)))
	}
	numberOfWorkerBuffers := numberOfWorkers * 4
	if numberOfWorkerBuffers > uint(len(input.PathList)) { numberOfWorkerBuffers = uint(len(input.PathList)) }
	list := make(chan string, numberOfWorkerBuffers)
	output := abstract.NewOutput(input.Stdout, input.Stderr)

	if input.Verbose {
		output.Write("using %v workers for %v inputs.\n", numberOfWorkers, len(input.PathList))
	}

	work := parallel.NewWork(numberOfWorkers)

	work.Feed(func() {
		for _, path := range input.PathList {
			list <-path
		}; close(list)
	})

	work.Start(func() {
		exitCode = hash1(input, output, factory, list)
	})

	work.Wait(); output.Done(); output.Wait()

	return
}

// hash one file, after that restart for another file read from list
func hash1(input *Input, output *abstract.Output, factory func() hash.Hash, list chan string) (exitCode uint8) {
	buffers := make(chan NamedBuffer, 50)
	work := parallel.NewWork(1); hasher := factory(); first := true

	// open and read in one thread
	work.Feed(func() {
		defer close(buffers)

		helper := iotool.ReadOnly().ToggleFileAdviceReadSequential()
		if input.NoCache { helper.ToggleFileAdviceDontNeed() }
		helper.SetE(func(name string, e error) {
			output.WriteEMessage(e, " on '%v'", name)
			if iotool.IsNotExist(e) {
				exitCode = abstract.ERROR_FILE_NOT_FOUND
			} else {
				exitCode = abstract.ERROR_READING
			}
		})

		ReadFiles(helper, buffers, list)
	})

	// hashing in another thread
	work.Run(func() {
		for buffered := range buffers {
			if buffered.reset && !first {
				hasher.Reset()
			} else {
				first = false
			}

			if buffered.done {
				output.Write("%v  %v\n", hex.EncodeToString(hasher.Sum(nil)), buffered.name)
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

// hash a file or stdin
func Md5(input *Input) (exitCode uint8) {
	return doHash(input, md5.New)
}

// hash a file or stdin
func Sha1(input *Input) (exitCode uint8) {
	return doHash(input, sha1.New)
}

// hash a file or stdin
func Sha224(input *Input) (exitCode uint8) {
	return doHash(input, sha256.New224)
}

// hash a file or stdin
func Sha256(input *Input) (exitCode uint8) {
	return doHash(input, sha256.New)
}

// hash a file or stdin
func Sha384(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New384)
}

// hash a file or stdin
func Sha512(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New)
}

// hash a file or stdin
func Sha512_224(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New512_224)
}

// hash a file or stdin
func Sha512_256(input *Input) (exitCode uint8) {
	return doHash(input, sha512.New512_256)
}

// hash a file or stdin
func Adler32(input *Input) (exitCode uint8) {
	return doHash(input, func() hash.Hash {
	    return adler32.New()
	})
}

// hash a file or stdin
func Crc32(input *Input) (exitCode uint8) {
	return doHash(input, func() hash.Hash {
		return crc32.NewIEEE()
	})
}

// hash a file or stdin
func Crc64(input *Input) (exitCode uint8) {
	return doHash(input, func() hash.Hash {
		return crc64.New(crc64.MakeTable(crc64.ISO))
	})
}

// hash a file or stdin
func Hash(input *Input) (exitCode uint8) {
	switch(input.Type) {
	case MD5:
		return Md5(input)
	case SHA1:
		return Sha1(input)
	case SHA224:
		return Sha224(input)
	case SHA256:
		return Sha256(input)
	case SHA384:
		return Sha384(input)
	case SHA512:
		return Sha512(input)
	case SHA512_224:
		return Sha512_224(input)
	case SHA512_256:
		return Sha512_256(input)
	case ADLER32:
		return Adler32(input)
	case CRC32:
		return Crc32(input)
	case CRC64:
		return Crc64(input)
	}

	return abstract.ERROR_NO_INPUT
}
