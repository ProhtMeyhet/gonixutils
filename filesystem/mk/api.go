package mk

import(

)

// create common directory structes (files, directorys and links)
func Mk(input *Input) (exitCode uint8) {
	switch {
	case input.Link:
		return Link(input)
	case input.Temporary:
		return Temporary(input)
	case input.File:
		return File(input)
	default:
		return Directory(input)
	}

}
