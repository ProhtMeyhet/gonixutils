package ls

import(
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func Ls(input *Input) (exitCode uint8) {
	work := parallel.NewStringFeeder(input.Paths)
	output := abstract.NewSortedTabbedOutput(input.Stdout, input.Stderr)
	if input.Lines {
		output = abstract.NewOutput(input.Stdout, input.Stderr)
	} else if input.NoSort {
		output = abstract.NewTabbedOutput(input.Stdout, input.Stderr)
	}

	if input.Detail {
		output.ToggleLinesManual()
		input.writeEntry = WriteEntryLong
	} else {
		input.writeEntry = WriteEntryShort
	}

	return list(input, output, work)
}

func WriteEntryLong(input *Input, output abstract.OutputInterface, entry *Entry) {
	name := input.decorate(entry.info.Name(), entry.info); owner := "?"; group := "???"
	mode := os.FileMode(0); modificationTime := time.Time{}; size := int64(0)
	if entry.info != nil {
		mode = entry.info.Mode(); modificationTime = entry.info.ModTime()
		size = entry.info.Size()

		sys := entry.info.Sys() // can be nil!
		if sys != nil {
			userId := strconv.Itoa(int(sys.(*syscall.Stat_t).Uid))
			groupId := int(sys.(*syscall.Stat_t).Gid)
			userinfo, e := user.LookupId(userId)
			if e != nil {
				owner = userId
			} else {
				owner = userinfo.Name
			}

			//FIXME in go 1.7 there will be user.LookupGroup
				group = strconv.Itoa(int(groupId))
		}
	}

	if mode & os.ModeSymlink == os.ModeSymlink {
		symlink, e := os.Readlink(entry.path)
		if e != nil {
			name += " -> broken link (" + e.Error() + ")"
		} else {
			name += " -> " + symlink
		}
	}

	output.WriteSorted("%v", name+"a", mode)
	output.WriteSorted("%v", name+"b", owner)
	output.WriteSorted("%v", name+"c", group)
	output.WriteSorted("%v", name+"d", size)
	output.WriteSorted("%v", name+"e", modificationTime.Format(time.Stamp))
	output.WriteSorted("%v\n", name+"f", name)
/*	output.Write("%v", mode)
	output.Write("%v", owner)
	output.Write("%v", group)
	output.Write("%v", size)
	output.Write("%v", modificationTime.Format(time.Stamp))
	output.Write("%v", name + "\n")*/
}

func WriteEntryShort(input *Input, output abstract.OutputInterface, entry *Entry) {
	format := "%v "; if input.Lines { format += "\n" }
	if entry.info != nil {
		decorated := input.decorate(entry.info.Name(), entry.info)
		output.Write(format, decorated)
	} else {
		decorated := input.decorate(entry.path, nil)
		output.Write(format, decorated)
	}
}

