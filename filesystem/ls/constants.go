package ls

import(
	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

const(
	// FIXME lookup standard decorates
	DIRECTORY	= "/"
	TEMPORARY	= "?"
	SYMLINK		= "@"
	DEVICE		= "!"
	NAMED_PIPE	= "|" // FIFO
	SOCKET		= "="
	CHARACTER_DEVICE = "^"
	EXECUTEABLE	= "*"
	// if it gets sticky, a german word is required
	STICKY		= "§"

	COLOR_SYMLINK		= abstract.TERMINAL_COLOR_CYAN
	COLOR_DIRECTORY		= abstract.TERMINAL_COLOR_BLUE
	COLOR_EXECUTEABLE	= abstract.TERMINAL_COLOR_GREEN
	COLOR_RESET		= abstract.TERMINAL_COLOR_RESET

	BACKGROUND_TEMPORARY	= abstract.TERMINAL_BACKGROUND_GREEN
)
