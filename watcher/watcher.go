package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/v4run/bob/blogger"
	"github.com/v4run/bob/builder"
	"github.com/v4run/bob/runner"
)

const (
	SLEEP_TIME = 400 // Time between successive checks for modification
)

/**
 * Watcher struct
 * b   : builder
 * r   : runner
 * dir : directory to watch
 */
type Watcher struct {
	b         builder.Builder
	r         runner.Runner
	dir       string
	buildOnly bool
}

/**
 * Returns a new watcher.
 */
func NewWatcher(path, appName, commands string, buildOnly bool) Watcher {
	if appName == "" {
		appName = filepath.Base(path)
	}
	return Watcher{dir: path, b: builder.NewBuilder(appName, path), r: runner.NewRunner(appName, path, commands), buildOnly: buildOnly}
}

/**
 * Called for each file in the directory.
 * If a change is found, a build is performed and the service is run.
 * Directories and files starting with `.` are skipped.  Returns error if any.
 */
func (w *Watcher) watchFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if strings.HasPrefix(filepath.Base(path), ".") { // skip directories like .git, .idea etc.
		if info.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	}

	if filepath.Ext(path) == ".go" {
		if info.ModTime().After(w.b.LastBuild()) {
			p, _ := filepath.Rel(w.dir, path)
			blogger.Info().Command("modified").Message(blogger.FormattedMessage(p)).Log()
			okay := w.b.Build()
			if okay {
				if !w.buildOnly {
					re := w.r.Run()
					return re
				}
				ce := w.r.RunCustom()
				if ce != nil {
					blogger.Error().Message(ce.Error()).Log()
				}
			}
		}
	}
	return nil
}

/**
 * Watch watches the path for changes.
 * An initial build and run is performed before watching begins. Returns error if any.
 */
func (w *Watcher) Watch() error {
	blogger.Info().Command("watching").Message(blogger.FormattedMessage(w.dir)).Log()

	// Do a first build
	okay := w.b.Build()
	if okay {
		if !w.buildOnly {
			// Do a first run.
			if re := w.r.Run(); re != nil {
				blogger.Error().Message(re.Error()).Log()
			}
		}
		ce := w.r.RunCustom()
		if ce != nil {
			blogger.Error().Message(ce.Error()).Log()
		}
	}
	stopWatch := make(chan error)
	go func() {
		ticker := time.NewTicker(time.Millisecond * SLEEP_TIME)
		for {
			err := filepath.Walk(w.dir, w.watchFunc)
			if err != nil && err != filepath.SkipDir {
				stopWatch <- err
				break
			}
			<-ticker.C
		}
	}()
	return <-stopWatch
}
