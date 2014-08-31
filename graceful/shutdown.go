package graceful

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/v4run/bob/bLogger"
)


/**
 * Handles graceful shutdown.
 */
func ActivateGracefulShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for sig := range signalChan {
		if sig == syscall.SIGINT {
			bLogger.Logger().Warn("Interrupt signal received. Exiting.")
			os.Exit(0)
		}
	}
}
