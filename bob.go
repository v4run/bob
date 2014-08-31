package main

import (
	"flag"
	"fmt"
	"github.com/v4run/bob/b_logger"
	"github.com/v4run/bob/graceful"
	"github.com/v4run/bob/watcher"
	"os"
	"path/filepath"
	"strings"
)

const (
	VERSION = "1.0.0"
)

var (
	path    string
	version bool
	help    bool
)

func init() {
	flag.StringVar(&path, "d", "", "-d <directory to watch>")
	flag.StringVar(&path, "dir", "", "-dir <directory to watch>")
	flag.BoolVar(&version, "v", false, "-v <directory to watch>")
	flag.BoolVar(&version, "version", false, "-version <directory to watch>")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bob [options]\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		fmt.Fprintf(os.Stderr, "\t-d, -dir       Directory   The directory to watch.\n")
		fmt.Fprintf(os.Stderr, "\t-v, -version   Version     Prints the version.\n")
		fmt.Fprintf(os.Stderr, "\t-h, -help      Help        Show this help.\n")
	}
}

func parseFlags() {
	flag.Parse()

	if version {
		fmt.Printf("bob v%s\n", VERSION)
		os.Exit(0)
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if strings.TrimSpace(path) == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		path = dir
	} else {
		dir, err := filepath.Abs(path)
		if err != nil {
			b_logger.Logger().Error("Invalid directory,", path)
			os.Exit(1)
		}
		path = dir
	}
}

func main() {
	go graceful.ActivateGracefulShutdown()
	parseFlags()
	w := watcher.NewWatcher(path)
	if err := w.Watch(); err != nil {
		b_logger.Logger().Error(err.Error())
		os.Exit(1)
	}
}
