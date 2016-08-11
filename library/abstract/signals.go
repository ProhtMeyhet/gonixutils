package abstract

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals() {
	go handleSignals()
}

func handleSignals() {
	signalChannel := make(chan os.Signal, 1)
        signal.Notify(signalChannel, syscall.SIGTERM)
        signal.Notify(signalChannel, syscall.SIGABRT)
        signal.Notify(signalChannel, syscall.SIGINT)

infinite:
	for {
		select {
		case <-signalChannel:
			fmt.Println()
			os.Exit(ERROR_COMPAT)
			break infinite
		}
	}
}
