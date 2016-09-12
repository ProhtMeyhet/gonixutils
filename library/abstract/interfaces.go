package abstract

import(
	"io"
	"sync"
)

type OutputInterface interface {
	sync.Locker
	io.Writer

	Initialise(io, e io.Writer)

	WriteFormatted(format string, values ...interface{})
	// buffer output if aplicable and sort it before outputting
	WriteSorted(format, sortkey string, values ...interface{})

	// used for sorted output, otherwise the same as Write
	Append(format string, values ...interface{})

	WriteE(e error) bool
	WriteEMessage(e error, format string, values ...interface{}) bool
	WriteError(format string, values ...interface{})

	// state you are done. close channels
	Done()

	// wait till everyting is written
	Wait()

	SortKey() int
	SetSortKey(to int)
	SetSortTransformator(to func(string) string)
	SortReversed() bool
	ToggleSortReversed()

	LinesManual() bool
	ToggleLinesManual()

	NewSubBuffer(name string, key int) OutputInterface
	PrintSubBufferNames() bool
	TogglePrintSubBufferNames()
	OrderBySubBufferNames() bool
	ToggleOrderBySubBufferNames()
	OrderBySubBufferKeys() bool
	ToggleOrderBySubBufferKeys()

	// internal
	writeTo(io.Writer)
}
