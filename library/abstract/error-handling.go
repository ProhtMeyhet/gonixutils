package abstract

import(
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ProhtMeyhet/libgosimpleton/simpleton"
)

func PrintError(e error, writer io.Writer, message string, extra ...interface{}) bool {
	return printError(e, false, writer, message, extra...)
}

func PrintErrorWithError(e error, writer io.Writer, message string, extra ...interface{}) bool {
	return printError(e, true, writer, message, extra...)
}

// if e != nil, exit
func ExitOnError(e error, writer io.Writer, exitCode uint8, message string, extra ...interface{}) {
	if e == nil { return }

	printError(e, false, writer, message, extra...)

	os.Exit(int(exitCode))
}

func printError(e error, printError bool, writer io.Writer, message string, extra ...interface{}) bool {
	if e == nil { return false }

	message = "%v " + message

	if len(extra) > 1 {
		message += strings.Repeat(" , '%v' ", len(extra) - 1) + "\n"
	}

	if simpleton.GetLastChar(message) != "\n" {
		message += "\n"
	}

	var tempExtra []interface{}
	tempExtra = append(tempExtra, os.Args[0])
	if printError {
		message = "%v " + message
		tempExtra = append(tempExtra, e.Error())
		tempExtra = append(tempExtra, extra...)
	} else {
		tempExtra = append(tempExtra, extra...)
	}

	fmt.Fprintf(writer, message, tempExtra...)

	return true
}
