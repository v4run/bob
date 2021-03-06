package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/v4run/bob/blogger"
	"github.com/v4run/bob/graceful"
	"github.com/v4run/bob/watcher"
)

const (
	VERSION = "1.0.1"
)

var (
	path        string
	envFilePath string
	name        string
	version     bool
	help        bool
	buildOnly   bool
	commands    string
)

func init() {
	flag.StringVar(&path, "p", "", "")
	flag.StringVar(&path, "path", "", "")
	flag.StringVar(&name, "n", "", "")
	flag.StringVar(&name, "name", "", "")
	flag.StringVar(&envFilePath, "e", "", "")
	flag.StringVar(&envFilePath, "env", "", "")
	flag.StringVar(&commands, "c", "", "")
	flag.StringVar(&commands, "commands", "", "")
	flag.BoolVar(&buildOnly, "b", false, "")
	flag.BoolVar(&buildOnly, "buildonly", false, "")
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bob [options]\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		fmt.Fprintf(os.Stderr, "\t-p, -path      Directory                The directory to watch.\n")
		fmt.Fprintf(os.Stderr, "\t-n, -name      Name                     The name for binary file.\n")
		fmt.Fprintf(os.Stderr, "\t-e, -env       Environment file path    Path to file containing environment variables to be set for the service.\n")
		fmt.Fprintf(os.Stderr, "\t-b, -buildonly  Build only mode         Just do a build when a change is detected.\n")
		fmt.Fprintf(os.Stderr, "\t-c, -commands   Custom commands to run  Run custom command after build.\n")
		fmt.Fprintf(os.Stderr, "\t-v, -version   Version                  Prints the version.\n")
		fmt.Fprintf(os.Stderr, "\t-h, -help      Help                     Show this help.\n")
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
		dir, _ := os.Getwd() // Get the current working directory if no path is specified explicitly
		path, _ = filepath.Abs(dir)
	} else {
		dir, err := os.Stat(path)
		if err != nil {
			blogger.Error().Message("Cannot find path,", blogger.FormattedMessage(path)).Log()
			os.Exit(1)
		}
		if !dir.IsDir() {
			blogger.Error().Message(fmt.Sprintf("Invalid path, %s. Path must be directory.", blogger.FormattedMessage(path))).Log()
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
		blogger.Warn().Message("Configuration file path provided is invalid,", blogger.FormattedMessage(path)).Log()
	} else {
		path, _ = filepath.Abs(path)
		blogger.Info().Command("exporting").Message(blogger.FormattedMessage(path)).Log()
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
	go graceful.ActivateGracefulShutdown() // go routine for watching interrupt signals
	parseFlags()
	if envFilePath != "" {
		setEnvs(envFilePath)
	}
	w := watcher.NewWatcher(path, name, commands, buildOnly)
	if err := w.Watch(); err != nil {
		blogger.Error().Message(err.Error()).Log()
		os.Exit(1)
	}
}
