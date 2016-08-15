package ls

import(
	"os"
)

func ShouldPrint(input *Input, name string) bool {
	if input.All { return true }
	if !input.Hidden && IsHidden(name) { return false }
	return true
}

func IsHidden(name string) bool {
	// that's a philosophical question, eh?
	if len(name) == 0 { return false }
	if name[0:1] == "." { return true }
	return false
}

type Entry struct {
	path, base string
	info os.FileInfo
}
