package echo

import(
	"fmt"
	"strings"
)

func Echo(input *Input) (exitCode uint8) {
	out := input.Stdout; if input.PrintToStderr { out = input.Stderr }

	printer := func(output ...string) { fmt.Fprintln(out, StringListToInterface(output...)...) }
	if input.NoNewLine { printer = func(output ...string) { fmt.Fprint(out, StringListToInterface(output...)...) } }

	if len(input.Arguments) == 0 {
		if input.NoNewLine { return }
		printer(); return
	}

	output := input.Arguments

	if input.Escapes {
		output = make([]string, 1); output[0] = strings.Join(input.Arguments, " ")
		escaped := ParseTerminalEscape([]byte(output[0]))
		output[0] = string(escaped)
	}

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
