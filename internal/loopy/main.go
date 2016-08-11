package main

import(
	"fmt"

	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

// something that does nothing but wait for a signal to kill
// for testing and debugging
func main() {
	fmt.Printf("me pid is: %v\nwating for signal to print and SIGINT to quit!\n", os.Getpid())
	handleSignals()
}

func handleSignals() {
	signalChannel := make(chan os.Signal, 1)
        signal.Notify(signalChannel)
	printStack, refuseExit := 0, false

infinite:
	for {
		select {
		case signal := <-signalChannel:
			fmt.Printf("got signal %v. ", signal)

			switch signal {
			case syscall.SIGKILL:
				fmt.Println("murderer!")
			case syscall.SIGABRT:
				if printStack < 2 {
					printStack++
					fmt.Println("ignoring!")
				}
			case syscall.SIGTSTP:
				fmt.Println("waiting for SIGCONT!")
			case syscall.SIGTERM:
				fmt.Println("exiting only on SIGINT!")
			case syscall.SIGUSR1:
				refuseExit = true // too cosy!
				fmt.Println("it's nice and cosy here! i think i'll stay!")
			case syscall.SIGUSR2:
				if refuseExit {
					refuseExit = false // not so cosy anymore!
					fmt.Println("pretty pretty pretty please don't do what i think you are about to do!")
				} else {
					fmt.Println("it's been revoked!")
				}
			case syscall.SIGINT:
				if refuseExit { fmt.Println(); continue }
				fmt.Println("exiting.")
				os.Exit(0)
				break infinite
			default:
				fmt.Println()
			}

			if printStack == 2 {
				pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
				printStack = 0
			}
		}
	}
}
