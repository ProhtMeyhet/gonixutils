package echo

import(
	"fmt"
)

func Echo(input *Input) (exitCode uint8) {
	// where to write to?
	out := input.Stdout; if input.PrintToStderr { out = input.Stderr }

	// define a little printer wrapper
	printer := func(output ...string) { fmt.Fprintln(out, StringListToInterface(output...)...) }
	if input.NoNewLine { printer = func(output ...string) { fmt.Fprint(out, StringListToInterface(output...)...) } }

	// just a newline is requested ... except if input.NoNewLine is given of course
	if len(input.Arguments) == 0 { printer(); return }

	// parse escape sequences if requested
	output := input.Arguments; if input.Escapes {
		output = ParseTerminalEscape(input.Arguments...)
	}

	// and finally print
	printer(output...); return
}

// TODO move to libgosimpleton
// FIXME: fix in go
func StringListToInterface(in ...string) (out []interface{}) {
	out = make([]interface{}, len(in))
	for key, value := range in {
		out[key] = value
	}; return
}
