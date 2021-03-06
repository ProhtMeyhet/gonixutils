package rm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// remove filesystem entries given in the input.PathList
// FIXME: use abstract.OutputInterface
func Rm(input *Input) (exitCode uint8) {
	work := parallel.NewStringsFeeder(input.PathList...)

	work.Start(func() {
		for path := range work.Talk {
			// TODO more fine grained locking if Interactive (stats parallel)
			if input.Interactive { work.Lock() }
			e, _ := Remove(input, path)
			if e != nil {
				io.WriteString(input.Stderr, e.Error() + "\n")
				// FIXME more fine grained error
				exitCode = abstract.FAILED
			} else if input.VerboseLevel >= 2 {
				io.WriteString(input.Stdout, fmt.Sprintf("deleted '%v'\n", path))
			}
			if input.Interactive { work.Unlock() }
		}
	}); work.Wait()

	return
}

// TODO asynchronously remove if Recursive

// remove one filesystem entry.
// with input.Recursive, remove all filesystem entries within a directory and the directory.
// with input.Interactive, ask user before remove.
func Remove(input *Input, path string) (e error, allRemoved bool) {
	stat, e := os.Lstat(path); if e != nil { return }

	if stat.IsDir() {
		return removeDir(path, input)
	}

	if input.Interactive {
		if stat.Size() == 0 {
			if !YesOrNo(input, "rm: do you want to remove regular empty file '%s' ? (y/N) ", path) {
				return nil, false
			}
		} else if !YesOrNo(input, "rm: do you want to remove regular file '%s'\n(size: %v b; modification time: %v) ? (y/N) ",
					path, strconv.FormatInt(stat.Size(), 10), stat.ModTime().Format(time.UnixDate)) {
			return nil, false
		}
	}

	return os.Remove(path), true
}

// removes a directory, if input.Recurive also all it's contents.
func removeDir(path string, input *Input) (e error, allRemoved bool) {
	count := 0

	names, e := iotool.ListDirectory(path); if e != nil { return }
	if input.Interactive {
		if len(names) == 0 {
			fmt.Fprintf(input.Stdout, "directory '%s' empty.\n", path)
		} else if !YesOrNo(input, "rm: descend into directory '%s'? (y/N) ", path) {
			return nil, false
		}
	}

	// allow removing of empty directories without input.Recursive
	if len(names) == 0 {
		goto remove
	}

	// cannot remove due to recursive not set
	if !input.Recursive && len(names) > 0 {
		return errors.New(fmt.Sprintf("rm: '%s': directory not empty!", path)), false
	}

	// remove recursivly
	for _, name := range names {
		deletePath := path + string(os.PathSeparator) + name
		e, _ = Remove(input, deletePath)
		// FIXME report e to output
		if e == nil { count++ }
	}; if count == len(names) { allRemoved = true }

	// recursivly remove failed, can't remove parent
	if !allRemoved {
		fmt.Fprintf(input.Stdout, "cannot remove '%s', not empty\n", path)
		return
	}

remove:
	if input.Interactive && !YesOrNo(input, "rm: remove directory '%s'? (y/N) ", path) {
		return
	}

	return os.Remove(path), allRemoved
}

// return true on "y", "Y" and "yes"; otherwise false
func YesOrNo(input *Input, question string, replace ...string) bool {
	answer := ""

	fmt.Fprintf(input.Stdout, fmt.Sprintf(question, ToInterfaceList(replace...)...))
	_, e := fmt.Fscanln(input.Stdin, &answer); if e != nil { goto out }

	if answer == "y" || answer == "Y" ||  answer == "yes" {
		return true
	}

out:
	return false
}

// alias
func Delete(input *Input, path string) (e error, allRemoved bool) {
	return Remove(input, path)
}

// TODO move to libgosimpleton
func ToInterfaceList(input ...string) (result []interface{}) {
	result = make([]interface{}, len(input))
	for key, entry := range input {
		result[key] = entry
	}; return
}
