package ls

import(
	"os"
	"strings"

	"github.com/ProhtMeyhet/libgosimpleton/iotool"
	"github.com/ProhtMeyhet/libgosimpleton/parallel"

	"github.com/ProhtMeyhet/gonixutils/library/abstract"
)
// TODO
// split output of two or more directories. currently they are printed as if they were one
func list(input *Input, output abstract.OutputInterface, work *parallel.WorkString) (exitCode uint8) {
	initialise(input, output)

	work.Start(func() {
		iwork := parallel.NewWork(work.Workers())
		entries := make(chan *Entry, iwork.SuggestBufferSize(0)); entry := &Entry{}

		iwork.Feed(func() {
			var e error
			for path := range work.Talk {
				entry.info, e = os.Lstat(path); if output.WriteE(e) { continue }
				if entry.info.IsDir() {
					if input.All {
						// TODO can an error happen? maybe if path is / ?
						infoDotDot, _ := os.Lstat(".." + string(os.PathSeparator) + path)
						input.writeEntry(input, output, &Entry{ path: ".", base: ".", info: entry.info })
						input.writeEntry(input, output, &Entry{ path: "..", base: "..", info: infoDotDot })
					}
					list, e := iotool.ListDirectory(path)
					if e != nil { output.WriteE(e); entries <-entry }
					for _, ipath := range list {
						entry = &Entry{}
						entry.path = path + string(os.PathSeparator) + ipath; entry.base = ipath
						entries <-entry
					}
				} else {
					entry.path = path; entry.base = path
					entries <-entry
				}

				entry = &Entry{}
			}; close(entries)
		})

		iwork.Start(func() {
			var e error
			for entry := range entries {
				// first see if path should be printed, avoid syscall
				if !ShouldPrint(input, entry.base) { continue }
				if entry.info == nil { entry.info, e = os.Lstat(entry.path); if output.WriteE(e) { continue } }
				input.writeEntry(input, output, entry)
			}
		})

		iwork.Wait() // required if not using tabwriter ?: ; if !input.Lines { output.Append("%v", "\n") }
	})

	work.Wait(); output.Done(); output.Wait(); return
}

func initialise(input *Input, output abstract.OutputInterface) {
	if input.NoColor { input.decorate = PlainDecorator } else { input.decorate = ColorDecorator }

	if !input.NoSort {
		// before sorting remove . from files
		// also have to remove color junk in front to enable sorting...
		output.SetSortTransformator(func(name string) string {
			if len(name) == 0 ||
				(len(name) == 1 && name[0:1] == ".") ||
				(len(name) == 2 && name[0:2] == "..") { return name }
			if name[0:1] == "." { return name[1:] }
			if !input.NoColor && strings.HasPrefix(name, abstract.TERMINAL_COLOR_JUNK) {
				name = name[len(abstract.TERMINAL_COLOR_JUNK_LENGTH):]
			}
			return name
		})
	}
    }
