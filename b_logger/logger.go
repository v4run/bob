// +build !windows

package b_logger

import (
	"log"
	"os"
	"strings"
	"time"
)

type bLogger struct {
	*log.Logger
}

type Log struct {
	command string
	message string
	level   string
}

var (
	error_c   = "\033[38;5;196m"
	info_c    = "\033[38;5;266m"
	warn_c    = "\033[38;5;214m"
	flags_c   = "\033[38;5;8m"
	name_c    = "\033[38;5;157m"
	message_c = "\033[38;5;243m"
	command_c = "\033[38;5;248m"
	reset_c   = "\033[0m"
	bl        = bLogger{log.New(os.Stdout, name_c+"bob "+reset_c+flags_c, 0)}
)

func (l Log) Command(command string) Log {
	l.command = command
	return l
}

func (l Log) Message(message ...string) Log {
	l.message = strings.Join(message, " ")
	return l
}

func Error() Log {
	return Log{level: "error"}
}

func Info() Log {
	return Log{level: "info"}
}

func Warn() Log {
	return Log{level: "warn"}
}

func (l Log) Log() {
	out := "[" + time.Now().Format("15:04:05") + "] "
	if l.command != "" {
		out += command_c + l.command + reset_c + " "
	}
	switch l.level {
	case "error":
		out += error_c
		break
	case "info":
		out += info_c
		break
	case "warn":
		out += warn_c
		break
	}
	out += l.message + reset_c
	bl.Print(out)
}

func FormattedMessage(message string) string {
	return command_c + "'" + reset_c + message_c + message + reset_c + command_c + "'" + reset_c
}
