package rm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// remove filesystem entries given in the input.PathList
func Rm(input *Input) (e error) {
	for _, path := range input.PathList {
		e, _ = Remove(path, input); if e != nil { io.WriteString(input.Stderr, e.Error() + "\n") }
	}

	return
}

// remove one filesystem entry.
// with input.Recursive, remove all filesystem entries within a directory and the directory.
// with input.Interactive, ask user before remove.
func Remove(path string, input *Input) (e error, interactiveAll bool) {
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
func removeDir(path string, input *Input) (e error, all bool) {
	fd, e := os.Open(path); if e != nil { return }
	names, e := fd.Readdirnames(-1); if e != nil { return }
	fd.Close()

	if input.Interactive {
		if len(names) == 0 {
			fmt.Fprintf(input.Stdout, "directory '%s' empty.\n", path)
		} else if !YesOrNo(input, "rm: descend into directory '%s'? (y/N) ", path) {
			return nil, false
		}
	}

	// allow removing of empty directories without input.Recursive
	if !input.Recursive && len(names) > 0 {
		return errors.New(fmt.Sprintf("rm: '%s': directory not empty!", path)), false
	}

	allRemoved := false
	for _, name := range names {
		e, allRemoved = RemoveAll(path + string(os.PathSeparator) + name, input)
		if e != nil { return e, false }
	}

	if !allRemoved {
		if input.Interactive { fmt.Fprintf(input.Stdout, "cannot remove '%s', not empty\n", path) }
		return nil, false
	}

	if input.Interactive && !YesOrNo(input, "rm: remove directory '%s'? (y/N) ", path) {
		return
	}

	return os.Remove(path), true
}

// return true on "y", "Y" and "yes"; otherwise false
func YesOrNo(input *Input, question string, replace ...string) bool {
	answer := ""

	fmt.Fprintf(input.Stdout, fmt.Sprintf(question, ToInterfaceList(replace)...))
	_, e := fmt.Fscanln(input.Stdin, &answer); if e != nil { goto out }

	if answer == "y" || answer == "Y" ||  answer == "yes" {
		return true
	}

out:
	return false
}

// TODO move to libgosimpleton
func ToInterfaceList(input []string) (result []interface{}) {
	result = make([]interface{}, len(input))
	for key, entry := range input {
		result[key] = entry
	}

	return
}
