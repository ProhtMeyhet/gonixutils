package ls

import(
	"os"
)

// just returns the name
func NoneDecorator(name string, _ os.FileInfo) string {
	return name
}

// add some special characters if special name
func PlainDecorator(name string, info os.FileInfo) string {
	if info == nil {
		return name
	}

	mode := info.Mode()
	if mode & os.ModeSticky == os.ModeSticky {
		name += STICKY
	}

	if mode & os.ModeTemporary == os.ModeTemporary {
		return name + TEMPORARY
	} else if mode & os.ModeSymlink == os.ModeSymlink {
		return name + SYMLINK
	} else if mode & os.ModeDevice == os.ModeDevice {
		if mode & os.ModeCharDevice == os.ModeCharDevice {
			return name + CHARACTER_DEVICE
		}
		return name + DEVICE
	} else if mode & os.ModeNamedPipe == os.ModeNamedPipe {
		return name + NAMED_PIPE
	} else if mode & os.ModeSocket == os.ModeSocket {
		return name + SOCKET
	} else if mode & os.ModeDir == os.ModeDir {
		return name + DIRECTORY
	} else if mode & 0111 == 0111 { // EXEcuteABLE
		return name + EXECUTEABLE
	}; return name
}

// (...) Little by little I could begin to enjoy the unprecedented colors and plays of shapes
// that persisted behind my closed eyes. Kaleidoscopic, fantastic images surged in on me,
// alternating, variegated, opening and then closing themselves in circles and spirals,
// exploding in colored fountains, rearranging and hybridizing themselves in constant flux (...)
func ColorDecorator(name string, info os.FileInfo) string {
	if info == nil || len(name) == 0 {
		return name
	}

	mode := info.Mode()
	switch {
	case mode & os.ModeTemporary == os.ModeTemporary:
		return BACKGROUND_TEMPORARY + name + COLOR_RESET
	case mode & os.ModeSymlink == os.ModeSymlink:
		return COLOR_SYMLINK + name + COLOR_RESET
	case mode & os.ModeDevice == os.ModeDevice:
		if mode & os.ModeCharDevice == os.ModeCharDevice {
			return name //+ COLOR_RESET
		}
		return name //+ COLOR_RESET 
	case mode & os.ModeNamedPipe == os.ModeNamedPipe:
		return name //+ COLOR_RESET
	case mode & os.ModeSocket == os.ModeSocket:
		return name //+ COLOR_RESET
	case mode & os.ModeDir == os.ModeDir:
		return COLOR_DIRECTORY + name + COLOR_RESET
	case mode & 0111 == 0111: // EXEcuteABLE
		return COLOR_EXECUTEABLE + name + COLOR_RESET
	}; return name
}
