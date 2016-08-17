package ls

import(
	"os"
	"os/user"
	pathTools "path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)

func List(input *Input, paths []string, mainOutput abstract.OutputInterface) (exitCode uint8) {
	work := parallel.NewStringFeeder(paths); initialise(input, mainOutput); directoryCount := 0
	work.Start(func() {
		iwork := parallel.NewWork(work.Workers())
		entries := make(chan *Entry, iwork.SuggestBufferSize(0)); entry := &Entry{ output: mainOutput }

		iwork.Feed(func() {
			var e error
			for path := range work.Talk {
				if path == "" { mainOutput.WriteError("Path is empty!") }
				entry.info, e = os.Lstat(path); if mainOutput.WriteE(e) { continue }
				if entry.info.IsDir() {
					if !input.Union {
						entry.output = mainOutput.NewSubBuffer(path + ":\n", 0)
						if directoryCount < 2 && !mainOutput.PrintSubBufferNames() {
							work.Lock()
							directoryCount++; if directoryCount == 2 {
								if !mainOutput.PrintSubBufferNames() {
									mainOutput.TogglePrintSubBufferNames()
								}
								if !mainOutput.OrderBySubBufferNames() {
									mainOutput.ToggleOrderBySubBufferNames()
								}
							}
							work.Unlock()
						}
					}
					if input.All { writeDotDot(input, entry) }
					list, e := iotool.ListDirectory(path); if entry.output.WriteE(e) { continue }
					for _, ipath := range list {
						entry = &Entry{ output: entry.output }
						entry.path = pathTools.Join(path, ipath); entry.base = ipath
						entries <-entry
					}
				} else {
					entry.path = path; entry.base = path
					entries <-entry
				}

				entry = &Entry{ output: mainOutput }
			}; close(entries)
		})

		iwork.Start(func() {
			var e error
			for entry := range entries {
				// first see if path should be printed, avoid syscall
				if !ShouldPrint(input, entry.base) { continue }
				if entry.info == nil { entry.info, e = os.Lstat(entry.path); if entry.output.WriteE(e) { continue } }
				input.writeEntry(input, entry)
				// TODO stay in this goroutine or go list() ?
				if input.Recursive && entry.info.IsDir() {
					subOutput := mainOutput.NewSubBuffer("", 0)
					subOutput.TogglePrintSubBufferNames()
					subOutput.ToggleOrderBySubBufferNames()
					list(input, []string{ entry.path }, subOutput)
				}
			}
		})

		iwork.Wait() // required if not using tabwriter ?: ; if !input.Lines { output.Append("%v", "\n") }
	})

	work.Wait(); return
}

func initialise(input *Input, mainOutput abstract.OutputInterface) {
	if input.NoColor { input.decorate = PlainDecorator } else { input.decorate = ColorDecorator }

	if input.Detail {
		mainOutput.ToggleLinesManual()
		input.writeEntry = WriteEntryLong
	} else {
		input.writeEntry = WriteEntryShort
	}

	if input.SortReversed { output.ToggleSortReversed() }

	if !input.NoSort {
		// before sorting remove . from files
		// also have to remove color junk in front to enable sorting...
		mainOutput.SetSortTransformator(func(name string) string {
			if len(name) == 0 ||
				(len(name) == 1 && name[0:1] == ".") ||
				(len(name) == 2 && name[0:2] == "..") { return name }
			if !input.NoColor && strings.HasPrefix(name, abstract.TERMINAL_COLOR_JUNK) {
				name = name[len(abstract.TERMINAL_COLOR_JUNK_LENGTH):]
			}
			if name[0:1] == "." { return name[1:] }
			return name
		})
	}

	if input.Recursive {
		mainOutput.TogglePrintSubBufferNames()
		mainOutput.ToggleOrderBySubBufferNames()
	}
}

func WriteEntryLong(input *Input, entry *Entry) {
	name := input.decorate(entry.base, entry.info); owner := "?"; group := "???"
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

	entry.output.Lock()
	entry.output.WriteSorted("%v", name+"a", mode)
	entry.output.WriteSorted("%v", name+"b", owner)
	entry.output.WriteSorted("%v", name+"c", group)
	entry.output.WriteSorted("%v", name+"d", size)
	entry.output.WriteSorted("%v", name+"e", modificationTime.Format(time.Stamp))
	entry.output.WriteSorted("%v\n", name+"f", name)
	entry.output.Unlock()
}

func WriteEntryShort(input *Input, entry *Entry) {
	format := "%v "; if input.Lines { format += "\n" }
	if entry.info != nil {
		decorated := input.decorate(entry.base, entry.info)
		entry.output.Write(format, decorated)
	} else {
		decorated := input.decorate(entry.base, nil)
		entry.output.Write(format, decorated)
	}
}

func writeDotDot(input *Input, entry *Entry) {
	// TODO linux: if os.Lstat("/..") just returns /. so no error. BSD? 
	infoDotDot, _ := os.Lstat(pathTools.Join(entry.path, ".."))
	input.writeEntry(input, &Entry{ path: ".", base: ".",
					info: entry.info, output: entry.output })
	input.writeEntry(input, &Entry{ path: "..", base: "..",
					info: infoDotDot, output: entry.output })
}

// pass Entries around
type Entry struct {
	path, base string
	info os.FileInfo
	output abstract.OutputInterface
}
