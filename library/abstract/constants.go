package abstract

const(
	// to debug or not to debug, that is the question
	DEBUG			= false

	// this token stands for stdin
	STDIN_TOKEN		= "-"

	// part file extension
	EXTENSION_PART		= ".part"

	// global setting for FADVICE_DONTNEED before reading
	SET_FILE_ADVICE_DONTNEED = true

	// TODO add more colors and colours
	TERMINAL_COLOR_JUNK	= "\x1b["
	TERMINAL_COLOR_JUNK_LENGTH = "\x1b[xxxxx"
	TERMINAL_COLOR_CYAN	= "\x1b[36;1m"
	TERMINAL_COLOR_BLUE	= "\x1b[34;1m"
	TERMINAL_COLOR_GREEN	= "\x1b[32;1m"
	TERMINAL_COLOR_RESET	= "\x1b[0m"

	// FIXME unknown escape sequence: e
	TERMINAL_BACKGROUND_GREEN = "\\e[42mGreen"
)
