package sync

import(
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"
)

// Your job is being a professor and researcher: That's one hell of a good excuse for some of the
// brain-damages of minix.
//
//  -- Linus Torvalds to Andrew Tanenbaum

// Synchronize all filesystems or a file.
func Sync(input *Input) (exitCode uint8) {
        output := abstract.NewOutput(input.Stdout, input.Stderr)
        exitCode = Synchronize(input, output)
        output.Done(); output.Wait(); return
}

// Synchronize all filesystems or a file.
//
// According to the standard specification (e.g., POSIX.1-2001), sync()
// schedules the writes, but may return before the actual writing is
// done.  However Linux waits for I/O completions, and thus sync() or
// syncfs() provide the same guarantees as fsync called on every file in
// the system or filesystem respectively.
func Synchronize(input *Input, output abstract.OutputInterface) (exitCode uint8) {
	mutuallyExclusive, e := MutuallyExclusiveBoolError(ERROR_MUTUALLY_EXCLUSIVE_OPTIONS,
								input.Data,
								input.File,
								input.FileSystem)
	if mutuallyExclusive {
		output.WriteE(e); exitCode = abstract.ERROR_INVALID_ARGUMENT; return
	}

	// commit all filesystem caches to disk. cannot fail
	if len(input.PathList) == 0 {
		syscall.Sync(); return
	}

	// TODO
	// fifo's (pipes) -> open with NONBLOCK. after that, discard NONBLOCK with fcntl
	// @see coreutils
	helper := iotool.WriteOnly().ToggleDoNotTestForDirectory()

	work := parallel.NewStringsFeeder(input.PathList...)
	work.Start(func() {
		for path := range work.Talk {
			file, e := iotool.Open(helper, path); if output.WriteE(e) {
				if exitCode == 0 { exitCode = abstract.ERROR_OPENING_FILE }
				continue
			}

			switch {
			case input.File:
				e = file.Sync()
			case input.FileSystem:
				e = unix.Syncfs(int(file.Fd()))
			case input.Data:
				e = syscall.Fdatasync(int(file.Fd()))
			default:
				e = file.Sync()
				if e == nil {
					directoryPath := filepath.Dir(path)
					directory, ie := iotool.OpenDirectory(directoryPath)
					if ie == nil {
						directory.Sync(); directory.Close()
					} else {
						e = ie
					}
				}
			}

			if output.WriteE(e) && exitCode == 0 { exitCode = abstract.PARTLY }

			file.Close(); file = nil
		}
	})

	work.Wait(); return
}

// returns true and the given error, if more then one value is true
func MutuallyExclusiveBoolError(inputE error, values ...bool) (bool, error) {
	isTrue := false
	for _, value := range values {
		if isTrue && value { return true, inputE }
		if value { isTrue = true }
	}

	return false, nil
}
