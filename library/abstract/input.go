package abstract

import(
	"io"
	"os"
)

type Input struct {
	Stdin	io.Reader
	Stdout	io.Writer
	Stderr	io.Writer

	VerboseLevel	uint

	Verbose, Help, Version	bool
}

func (input *Input) InitCli() {
	input.Stdin = os.Stdin
	input.Stdout = os.Stdout
	input.Stderr = os.Stderr
}
