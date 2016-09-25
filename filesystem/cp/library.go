package cp

import(
	"errors"
	"io"
	"os"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// TODO --recursive
// TODO --synchronize

func Cp(input *Input) (exitCode uint8) {
	if len(input.Paths) == 0 {
		return abstract.ERROR_NO_INPUT
	}

	output := abstract.NewOutput(input.Stdout, input.Stderr)

	if input.Destination == "" && len(input.Paths) == 1 { var e error
		input.Destination, e = os.Getwd(); if output.WriteE(e) { return abstract.ERROR_UNHANDLED }
	}

	helper := prepareFileHelper(input, output, &exitCode)

	exitCodeCopy := Copy(input, output, helper); if exitCodeCopy != 0 { exitCode = exitCodeCopy }

	output.Done(); output.Wait(); return
}

func Copy(input *Input, output abstract.OutputInterface, helper *iotool.FileHelper) (exitCode uint8) {
	copyChannel := make(chan *iotool.Cp, parallel.SuggestBufferSize(uint(len(input.Paths)))); cp := &iotool.Cp{}; var e error
	writeHelper := iotool.WriteOnly().ToggleCreate().ToggleExclusive()

	work := CopyFilesystemEntriesCreateDone(copyChannel, input.Recursive, func(cp *iotool.Cp) (handler iotool.FileInterface) {
		e = checkFileInfo(input, writeHelper, cp.From, cp.To); if e != nil {
			helper.RaiseError(cp.From, e); return nil
		}
		handler, e = createDestination(input, writeHelper, cp.From, cp.To); if e != nil {
			helper.RaiseError(cp.From, e); return nil
		}; return
	}, func(cp *iotool.Cp, _ int64) {
		// TODO
		/*if input.Synchronize && written != fileInfo.Size() {
			message := "short write on '" + cp.From + "'"
			if input.RemoveSource { message += " not removing source" }
			if input.RemoveSource && !input.NoPartFile { message += " and" }
			if !input.NoPartFile { message += " not renaming .part file" }
			e = errors.New(message); return
		}*/
		if !input.NoPartFile {
			e = os.Rename(cp.To + abstract.EXTENSION_PART, cp.To)
			if e != nil { return } // TODO write a good error message
		}
		if input.RemoveSource { e = os.Remove(cp.From); if e != nil { return } }
	})

	// FIXME somehow cache the FileInfo
	info, e := helper.FileInfo(input.Destination, true); if e == nil && info.IsDir() {
		for _, path := range input.Paths {
			cp.From = path; cp.To = PathAddPrefix(input.Destination, path)
			copyChannel <-cp
			cp = &iotool.Cp{}
		}
	} else if len(input.Paths) > 1 {
		output.WriteError("target '%v' is not a directory\n", input.Destination)
		return ERROR_DESTINATION_NOT_DIRECTORY
	} else if e == nil && !input.Force {
		output.WriteError("file '%v' exist. not overwriting without --force flag (1)\n", input.Destination)
		return abstract.ERROR_FILE_EXIST
	} else if !iotool.IsNotExist(e) && output.WriteE(e) {
		return abstract.ERROR_STAT
	} else {
		cp.From = input.Paths[0]; cp.To = input.Destination
		copyChannel <-cp
	}

	close(copyChannel); work.Wait(); return
}

func CopyFilesystemEntriesCreateDone(copyChannel chan *iotool.Cp, recursive bool,
		create func(cp *iotool.Cp) iotool.FileInterface,
		done func(cp *iotool.Cp, written int64)) (work parallel.WorkInterface) {
	// FIXME
	// if create == nil { create = CopyCreate }
	work = parallel.NewWork(0); helper := iotool.ReadOnly().ToggleFileAdviceReadSequential()

	work.Start(func() {
		var handler iotool.FileInterface; written := int64(0); var e error;
		for cp := range copyChannel {
			parallel.OpenFileDoWork(helper, cp.From, func(buffered *iotool.NamedBuffer) {
				handler = create(cp)
				if handler == nil { buffered.Cancel(); return }; defer handler.Close()
				written, e = io.Copy(handler, buffered); if e != nil { buffered.Cancel(); }
				if done != nil { done(cp, written) }
			}).Wait()
		}
	}); return
}

func createDestination(input *Input, writeHelper *iotool.FileHelper, from, destination string) (handler iotool.FileInterface, e error) {
	part := abstract.EXTENSION_PART; if input.NoPartFile { part = "" }

	if !input.NoPartFile {
		destination += part
	}

again:
	handler, e = iotool.Open(writeHelper, destination); if e != nil {
		if iotool.IsExist(e) {
			if !input.NoPartFile {
				if !input.DoNotRemovePartFile {
					writeHelper = iotool.WriteOnly().ToggleTruncate()
					goto again
				} else {
					e = errors.New("not removing part file: " + destination)
				}
			} else if input.Force {
				writeHelper = iotool.WriteOnly().ToggleTruncate()
				goto again
			} else {
				e = errors.New("file '" + destination +
					"' exist. not overwriting without --force flag . (2)")
			}
		}

		// if e is still nil there isn't anything that can be done here
		if e != nil { return }
	}

	e = applyPermissions(input, writeHelper, destination, from)

	return
}

func checkFileInfo(input *Input, helper *iotool.FileHelper, sourcePath, destinationPath string) (e error) {
	// if source file can't be stat'ed game over
	sourceInfo, e := helper.FileInfo(sourcePath, false); if e != nil { return }
	destinationInfo, e := helper.FileInfo(destinationPath, false); if e != nil {
		if iotool.IsNotExist(e) {
			e = nil
		// all other errors are considered in the sake of posix as not accessible
		} else if input.ForceRemoveIfNotAccessible {
			e = os.Remove(destinationPath)
		}

		// if the Remove has worked, e should be nil
		if e != nil { return }
	}

	if sourceInfo.Same(destinationInfo) {
		return errors.New("'" + sourcePath + "' and '" + destinationPath + "' are the same file")
	}

	return
}

func applyPermissions(input *Input, helper *iotool.FileHelper, destinationPath, sourcePath string) (e error) {
	sourceInfo, e := helper.FileInfo(sourcePath, true); if e != nil { return }
	permissions := sourceInfo.Mode()
	if !input.PreserveExecute && sourceInfo.IsExecuteable() { permissions ^= 0111 }
	e = os.Chmod(destinationPath, permissions)
	if input.PreserveOwner { e = os.Chown(destinationPath, sourceInfo.UserId(), sourceInfo.GroupId()) }

	return
}

/*
func CopyFilesystemEntries(from, to []string, recursive bool) (e error) {
	return
}

func CopyFilesystemEntriesCreate(from, to []string, recursive bool, create func(*iotool.NamedBuffer, string) iotool.FileInterface) (e error) {
	return
}

func CopyFilesystemEntriesDone(from, to []string, recursive bool, done func(*iotool.NamedBuffer, bool)) (e error) {
	return
}*/
