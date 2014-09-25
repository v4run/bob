package graceful

import (
	"github.com/v4run/bob/b_logger"
	"os"
	"os/signal"
	"syscall"
)

/**
 * Handles graceful shutdown.
 */
func ActivateGracefulShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for sig := range signalChan {
		if sig == syscall.SIGINT {
			b_logger.Warn().Command("interrupt").Message("Exiting.").Log()
			os.Exit(0)
		}
	}
}
