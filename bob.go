package main

import (
	"flag"
	"fmt"
	"github.com/v4run/bob/b_logger"
	"github.com/v4run/bob/graceful"
	"github.com/v4run/bob/watcher"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

const (
	VERSION = "1.0.0"
)

var (
	path        string
	envFilePath string
	name        string
	version     bool
	help        bool
)

func init() {
	flag.StringVar(&path, "p", "", "")
	flag.StringVar(&path, "path", "", "")
	flag.StringVar(&name, "n", "", "")
	flag.StringVar(&name, "name", "", "")
	flag.StringVar(&envFilePath, "e", "", "")
	flag.StringVar(&envFilePath, "env", "", "")
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bob [options]\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		fmt.Fprintf(os.Stderr, "\t-p, -path      Directory              The directory to watch.\n")
		fmt.Fprintf(os.Stderr, "\t-n, -name      Name                   The name for binary file.\n")
		fmt.Fprintf(os.Stderr, "\t-e, -env       Environment file path  Path to file containing environment variables to be set for the service.\n")
		fmt.Fprintf(os.Stderr, "\t-v, -version   Version                Prints the version.\n")
		fmt.Fprintf(os.Stderr, "\t-h, -help      Help                   Show this help.\n")
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
	validateFlags()
}

/**
 * Validate the flag values.
 */

func validateFlags() {
	if path == "" {
		path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	} else {
		dir, err := os.Stat(path)
		if err != nil {
			b_logger.Logger().Error("Cannot find path,", path)
			os.Exit(1)
		}
		if !dir.IsDir() {
			b_logger.Logger().Error(fmt.Sprintf("Invalid path, %s. Path must be directory.", path))
			os.Exit(1)
		}
		path, _ = filepath.Abs(path)
	}
}

/**
 * Sets the environment variables required for the service from configuration file.
 */
func setEnvs(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		b_logger.Logger().Warn("Configuration file path provided is invalid,", path)
	} else {
		values := strings.Split(string(contents), "\n")
		for _, value := range values {
			if strings.Contains(value, "=") {
				envVar := strings.Split(value, "=")
				os.Setenv(envVar[0], envVar[1])
			}
		}

	}
}

func main() {
	go graceful.ActivateGracefulShutdown()
	parseFlags()
	if envFilePath != "" {
		setEnvs(envFilePath)
	}
	w := watcher.NewWatcher(path, name)
	if err := w.Watch(); err != nil {
		b_logger.Logger().Error(err.Error())
		os.Exit(1)
	}
}
