package main

import (
	"github.com/v4run/bob/bLogger"
	"github.com/v4run/bob/watcher"
	"os"
	"path/filepath"
	"flag"
	"fmt"
	"github.com/v4run/bob/graceful"
)

const (
	VERSION = "1.0.0"
)

var (
	path string
	version bool
	help bool
	l = bLogger.Logger()
)

func init() {
	flag.StringVar(&path, "d", "", "-d <directory to watch>")
	flag.StringVar(&path, "dir", "", "-dir <directory to watch>")
	flag.BoolVar(&version, "v", false, "-v <directory to watch>")
	flag.BoolVar(&version, "version", false, "-version <directory to watch>")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bob [options]\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		fmt.Fprintf(os.Stderr, "\t-d, -directory Directory   The directory to watch.\n")
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

	if path == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		path = dir
	}
}

func main() {
	go graceful.ActivateGracefulShutdown()
	parseFlags()
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	w := watcher.NewWatcher(dir)
	if err := w.Watch(); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
}
