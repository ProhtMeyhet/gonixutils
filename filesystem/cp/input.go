package cp

import(

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

type Input struct {
	abstract.Input

	Paths		[]string
	Destination	string

	// unix mv
	RemoveSource	bool

	// do not preserve file attributes
	NoPreserve	bool

	// also preserve execute bit
	PreserveExecute	bool

	// also preserve owner
	PreserveOwner	bool

	// dont copy to .part file and rename, copy directly
	NoPartFile	bool

	// do not remove part file if it exists and do no copy
	DoNotRemovePartFile bool

	// overwrite
	Force		bool

	// posix; remove if stat failed
	// If a file descriptor for a destination file cannot be obtained attempt to unlink the destination file and proceed
	// should read: if a destination can't be accessed or written to, try to remove it
	ForceRemoveIfNotAccessible bool

	// copy recursive
	Recursive	bool

	// ask before overwrite
	Interactive	bool

	// TODO don't continue copying if possible
	NoContinue	bool

	// TODO call fsync on handler when done
	Sync		bool


	// TODO try to do a reflink on Cow filesystems like btrfs
	Reflink		bool

	NoCache		bool
}

// get an input for unix mv
func Mv() (input *Input) {
	input.RemoveSource = true; return
}

// get an input with posix compatible options set
func Posix() (input *Input) {
	input.Posix(); return
}

// get an input with strictly posix compatible options and boring set
func Strict() (input *Input) {
	input.Strict(); return
}

// set posix compatible options.
func (input *Input) Posix() {
	input.NoPreserve			= true
	input.Force				= true
	input.ForceRemoveIfNotAccessible	= true
	input.PreserveExecute			= true
	input.PreserveOwner			= false
}

// be boring
func (input *Input) Strict() {
	input.Posix()
	input.NoPartFile			= true
	input.NoContinue			= true
}
