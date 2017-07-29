package sync

import(
	"errors"
)

var ERROR_BAD_FILE_DESCRIPTOR		= errors.New("bad file descriptor!")
var ERROR_MUTUALLY_EXCLUSIVE_OPTIONS	= errors.New("options data, file & filesystem are mutually exclusive!")
