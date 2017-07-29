package sync

import(
	"syscall"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
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
	for _, path := range input.PathList {
		file, e := iotool.Open(helper, path); if output.WriteE(e) {
			if exitCode == 0 { exitCode = abstract.ERROR_OPENING_FILE }
			continue
		}

		switch {
		case input.File:
			e = syscall.Fsync(int(file.Fd()))
		case input.FileSystem:
			output.WriteError("Filesystem sync of [FILE] is not implemented!")
		/*
			// TODO, SYNCFS is not in syscall for x86 & x86_64
			_, _, errNo := syscall.Syscall(syscall.SYS_SYNCFS, file.Fd(), 0, 0, 0)
			if errNo > 0 {
				if errNo == syscall.EABDF {
					e = ERROR_BAD_FILE_DESCRIPTOR
				}
			}
		*/
		default:
			e = syscall.Fdatasync(int(file.Fd()))
		}

		if output.WriteE(e) && exitCode == 0 { exitCode = abstract.PARTLY }

		file.Close(); file = nil
	}

	return
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
