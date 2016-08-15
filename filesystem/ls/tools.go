package ls

import(
	"path"
)

// return true if the name should be printed
func ShouldPrint(input *Input, name string) bool {
	if input.All { return true }
	name = path.Base(name)
	if !input.Hidden && IsHidden(name) { return false }
	return true
}

// return true if the file is a unix hidden file
func IsHidden(name string) bool {
	// that's a philosophical question, eh?
	if len(name) == 0 { return false }
	if name[0:1] == "." { return true }
	return false
}
