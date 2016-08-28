package echo

import(

)

func ParseTerminalEscape(input ...string) (output []string) {
	output = make([]string, len(input))
	for key := range input {
		output[key] = ParseTerminalEscape1(input[key])
	}; return
}

func ParseTerminalEscape1(input string) string {
	return string(ParseTerminalEscapeBytes([]byte(input)))
}

func ParseTerminalEscapeBytes(input []byte) []byte {
	ii := 0
scan:
	for i := 0; i < len(input); {
		c := input[i]; i++
		if c == '\\' && i < len(input) {
			c = input[i]; i++
                        switch c {
			case 'a':
				c = '\a'
			case 'b':
				c = '\b'
			case 'c':
				break scan
			case 'e':
				c = '\x1B'
			case 'f':
				c = '\f'
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			case 'v':
				c = '\v'
			case '\\':
				c = '\\'
			case 'x':
				c = input[i]; i++
				if '9' >= c && c >= '0' && i < len(input) {
					hex := (c - '0')
					c = input[i]; i++
					if '9' >= c && c >= '0' && i <= len(input) {
						c = 16 * (c-'0') + hex
					}
				}
			}
		}
		input[ii] = c; ii++
        }

	return input[:ii]
}
