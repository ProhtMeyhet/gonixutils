package version

import(
	"fmt"
)

const(
	// only changes with incompatible changes
	MAJOR_VERSION = "0"

	// new features
	MINOR_VERSION = "1"

	// no incompatible changes, just bug fixes
	BUGFIX_VERSION = "0"

	// munchies
	GONIXUTILS_VERSION = MAJOR_VERSION + "." + MINOR_VERSION + "." + BUGFIX_VERSION
)

// just print a line with version information
func Print() {
	fmt.Println("gonixutils version: " + GONIXUTILS_VERSION)
}
