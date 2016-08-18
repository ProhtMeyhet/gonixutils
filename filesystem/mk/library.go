package mk

import(
	"fmt"
	"path/filepath"
	"os"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

// create links, either symbolic or hard
func Link(input *Input) (exitCode uint8) {
	if len(input.PathList) == 0 { exitCode = abstract.ERROR_NO_INPUT; return }

	target := input.PathList[0]
	if len(input.PathList) == 1 {// variant 2
		// use the file name of the target as link name
		linkname := filepath.Base(input.PathList[0])

		var e error
		if input.Symbolic {
		    e = os.Symlink(target, linkname)
		} else {
		    e = os.Link(target, linkname)
		}

		if abstract.PrintErrorWithError(e, input.Stderr, "") {
			exitCode = ERROR_LINKING
		}

		return
	}

	// variant 1
	// user has supplied link name
	input.PathList = input.PathList[1:]
	for _, path := range input.PathList {
		linkname := filepath.Base(path)

		var e error
		if input.Symbolic {
		    e = os.Symlink(target, linkname)
		} else {
		    e = os.Link(target, linkname)
		}

		if abstract.PrintErrorWithError(e, input.Stderr, "") {
			exitCode = ERROR_LINKING
		}
	}

	return
}

// create temporary files with unique name under os.TempDir()
func Temporary(input *Input) (exitCode uint8) {
	helper := iotool.WriteOnly()

	// use an empty entry for no prefix
	if len(input.PathList) == 0 {
		input.PathList = append(input.PathList, "")
	}

	for _, path := range input.PathList {
		handler, e := iotool.Temporary(helper, path)
		if e != nil {
			abstract.PrintErrorWithError(e, input.Stderr, "")
			exitCode = ERROR_CREATE_FILE
		} else {
			handler.Close()
			fmt.Fprintf(input.Stdout, "%v\n", handler.Name())
		}
	}

	return
}

// create files
func File(input *Input) (exitCode uint8) {
	helper := iotool.WriteOnly().ToggleExclusive()

	for _, path := range input.PathList {
		if input.Recursive {
			directory := filepath.Dir(path)
			// TODO mask
			if e := os.MkdirAll(directory, 0766); e != nil {
				abstract.PrintErrorWithError(e, input.Stderr, "")
				continue
			}
		}

		handler, e := iotool.Create(helper, path)

		if e != nil {
			abstract.PrintError(e, input.Stderr, "")
			exitCode = ERROR_CREATE_FILE
		} else {
			if input.Verbose {
				fmt.Fprintf(input.Stdout, "%v\n", path)
			}
			handler.Close()
		}
	}

	return
}

// create directorys
func Directory(input *Input) (exitCode uint8) {
	for _, path := range input.PathList {
		if input.Recursive {
			e := os.MkdirAll(path, 0766)
			if !abstract.PrintErrorWithError(e, input.Stderr, "") && input.Verbose {
				fmt.Fprintf(input.Stdout, "%v\n", path)
			}
		} else {
			e := os.Mkdir(path, 0766)
			if !abstract.PrintErrorWithError(e, input.Stderr, "") && input.Verbose {
				fmt.Fprintf(input.Stdout, "%v\n", path)
			}
		}
	}

	return
}
