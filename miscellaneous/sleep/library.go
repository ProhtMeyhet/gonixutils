package sleep

import(
	"fmt"

	"os"
	"os/signal"
	"time"
	"syscall"
)

var now time.Time

func Sleep(input *Input) (exitCode uint8) {
	duration := time.Duration(1); if now.IsZero() { now = time.Now() }
	if input.DoUntil {
	again:
		if !now.Before(input.Until) {
			fmt.Fprintf(input.Stderr, FluxCapacitorMalfunction(input.Until).Error())
			return EXIT_TIMETRAVEL_NOT_INVENTED_YET
		}

		duration = input.Until.Sub(now); if duration <= 0 { goto again }
	} else {
		duration = Sum(input.Durations...)
	}

	if input.Verbose {
		// FIXME drop the nanoseconds from duration output
		fmt.Fprintf(input.Stdout, "sleeping for %v (pid: %v)\n", duration, os.Getpid())
	}

	time.Sleep(duration)

	if input.VerboseLevel >= 2 {
		if input.Format == "" { input.Format = STAMP2 }
		fmt.Fprintf(input.Stdout, "It is now %v - Thank you for choosing gonixutils.\n", time.Now().Format(input.Format))
	}

	return
}

// go, catch signals
func SignalHandler(input *Input) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

again:
	select {
	case signal := <-signals:
		switch signal {
		case syscall.SIGINT, syscall.SIGTERM:
			fmt.Fprintf(input.Stderr, "sleep: catched signal %v after %v.\n", signal, time.Now().Sub(now))
			os.Exit(EXIT_INTERRUPTED)
		case syscall.SIGUSR1:
			fmt.Fprintf(input.Stderr, "slept for %v until now.\n", time.Now().Sub(now))
			goto again
		}
	}
}

// resets the globally used time.now
func Reset() {
	now = time.Now()
}
