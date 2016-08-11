package abstract

const(
	// try not to open more then this amount of goroutines for one task
	GOROUTINE_LIMIT			= 500
	// perform a runtime.GC() followed by runtime.Gosched() when hitting GOROUTINE_LIMIT
	GOROUTINE_LIMIT_GC		= true
	// try not to open more then this amount of goroutines for resources (files!)
	// must be divideable by 5 without leftovers
	GOROUTINE_LIMIT_RESOURCE	= 50
	// 1 of 5
	GOROUTINE_LIMIT_RESOURCE_1OF5	= GOROUTINE_LIMIT_RESOURCE / 5
	// 4 of 5
	GOROUTINE_LIMIT_RESOURCE_4OF5	= GOROUTINE_LIMIT_RESOURCE_1OF5 * 4

	// the smallest read for an io.Reader buffer
	READ_BUFFER_SMALL_SIZE		= 512
	// the usual read size for an io.Reader buffer
	READ_BUFFER_SIZE		= 4 * 1024
	// the maximum read size for an io.Reader buffer
	READ_BUFFER_MAXIMUM_SIZE	= 32 * 1024
)
