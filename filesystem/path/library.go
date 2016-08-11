package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// TODO:
// exitCodes testing and setting

func Path(input *Input) (exitCode uint8) {
	switch {
	case input.IsAbsolute:
		return IsAbsolute(input)
	case input.Clean:
		return Clean(input)
	case input.Extension:
		return Extension(input)
	case input.List:
		return List(input)
	case input.Basename:
		return Basename(input)
	case input.Directory:
		return Directory(input)
	case input.Join:
		return Join(input)
	case input.Split:
		return Split(input)
	case input.Match:
		return Match(input)
	case input.Relative:
		return Relative(input)
	default:
		return Absolute(input)
	}
}

func Split(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		splitted := strings.Split(path, string(os.PathSeparator))
		// TODO print empty line or not?
		// if len(splitted) > 0 && splitted[0] == "" { splitted = splitted[1:] }
		output(input, splitted...)
	}

	return
}

func Match(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		if matched, e := filepath.Match(input.Pattern, path); e != nil {
			abstract.PrintError(e, input.Stderr, "")
			exitCode = abstract.FAILED
		} else if !matched {
			exitCode = abstract.FAILED
		}
	}

	return
}

func Join(input *Input) (exitCode uint8) {
	output(input, filepath.Join(input.PathList...))
	return
}

func IsAbsolute(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		if !filepath.IsAbs(path) {
			return abstract.FAILED
		}
	}

	return
}

func Absolute(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		newpath, e := filepath.Abs(path); if e != nil {
			fmt.Fprintf(input.Stderr, "%v\n", e)
			exitCode = abstract.ERROR_INVALID_ARGUMENT
		} else {
			output(input, newpath)

			if exitCode > abstract.SUCCESSFUL && len(input.PathList) > 1 {
				exitCode = abstract.PARTLY
			}
		}
	}

	return
}

func Basename(input *Input) (exitCode uint8) {
	var e error

	for _, path := range input.PathList {
		if !input.Relative {
			path, e = filepath.Abs(path); if e != nil {
				fmt.Fprintf(input.Stderr, "%v\n", e)
				exitCode = abstract.ERROR_INVALID_ARGUMENT
				continue
			}
		}

		newpath := filepath.Base(path)
		output(input, newpath)

		if exitCode > abstract.SUCCESSFUL && len(input.PathList) > 1 {
			exitCode = abstract.PARTLY
		}
	}

	return
}

func Directory(input *Input) (exitCode uint8) {
	var e error

	for _, path := range input.PathList {
		if !input.Relative {
			path, e = filepath.Abs(path); if e != nil {
				fmt.Fprintf(input.Stderr, "%v\n", e)
				exitCode = abstract.ERROR_INVALID_ARGUMENT
				continue
			}
		}

		newpath := filepath.Dir(path)
		output(input, newpath)

		if exitCode > abstract.SUCCESSFUL && len(input.PathList) > 1 {
			exitCode = abstract.PARTLY
		}
	}

	return
}

func Extension(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		newpath := filepath.Ext(path)
		if newpath == "" {
			exitCode = abstract.EMPTY_OUTPUT
		} else {
			output(input, newpath)

			if exitCode == abstract.EMPTY_OUTPUT && len(input.PathList) > 1 {
				exitCode = abstract.PARTLY_EMPTY
			}
		}
	}

	return
}

func Relative(input *Input) (exitCode uint8) {
	if input.WorkingDirectory == "" {
		exitCode = abstract.ERROR_WORKING_DIRECTORY
		goto out
	}

	for _, path := range input.PathList {
		newpath, e := filepath.Rel(input.WorkingDirectory, path); if e != nil {
			if absolute, e := filepath.Abs(filepath.Dir(path)); e == nil && absolute == input.WorkingDirectory {
				output(input, path)
				continue
			}

			fmt.Fprintf(input.Stderr, "%v\n", e)
			exitCode = abstract.ERROR_INVALID_ARGUMENT
		} else {
			output(input, newpath)

			if exitCode > abstract.SUCCESSFUL && len(input.PathList) > 1 {
				exitCode = abstract.PARTLY
			}
		}
	}

out:
	return
}

func Clean(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		output(input, filepath.Clean(path))
	}

	return
}

func List(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		output(input, filepath.SplitList(path)...)
	}

	return
}

func output(input *Input, newpath ...string) {
	for _, path := range newpath {
		if input.Prefix != "" {
			path = strings.TrimPrefix(path, input.Prefix)
		}

		if input.Suffix != "" {
			path = strings.TrimSuffix(path, input.Suffix)
		}

		fmt.Fprintf(input.Stdout, "%v\n", path)
	}
}
