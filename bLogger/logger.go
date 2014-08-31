package bLogger

import (
	"log"
	"os"
	"strings"
)

type bLogger struct {
	*log.Logger
}

var (
	error_c = "\033[38;5;196m"
	info_c  = "\033[38;5;87m"
	warn_c  = "\033[38;5;214m"
	name_c  = "\033[38;5;157m"
	flags_c = "\033[38;5;241m"
	reset_c = "\033[0m"
	bl      = bLogger{log.New(os.Stdout, "[" + name_c+ "bob" + reset_c + "] "+flags_c, log.LstdFlags)}
)

func Logger() bLogger {
	return bl
}

func (b bLogger) printMessage(level string, args ...string) {
	message := strings.Join(args, " ");
	b.Printf("%s%s%s%s", reset_c, level, message, reset_c)
}

func (b bLogger) Info(args ...string) {
	b.printMessage(info_c, args...)
}

func (b bLogger) Error(args ...string) {
	b.printMessage(error_c, args...)
}

func (b bLogger) Warn(args ...string) {
	b.printMessage(warn_c, args...)
}
