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
    bl = bLogger{log.New(os.Stdout, "bob ", 0)}
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
        out += l.command + " "
    }
    out += l.message
    bl.Print(out)
}

func FormattedMessage(message string) string {
    return "'" + message + "'"
}
