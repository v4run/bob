package watcher

import (
	"github.com/v4run/bob/b_logger"
	"github.com/v4run/bob/builder"
	"github.com/v4run/bob/runner"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SLEEP_TIME = 400 // Time between successive checks for modification
)

/**
 * Watcher struct
 * b: builder
 * r: runner
 * dir: directory to watch
 */
type Watcher struct {
	b      builder.Builder
	r      runner.Runner
	dir   string
}

/**
 * Returns a new watcher.
 */
func NewWatcher(path, appName string) Watcher {
	if appName == "" {
		appName = filepath.Base(path)
	}
	return Watcher{dir: path, b: builder.NewBuilder(appName, path), r: runner.NewRunner(appName, path)}
}

func (w *Watcher) watchFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if strings.HasPrefix(filepath.Base(path), ".") { // skip directories like .git, .idea etc.
		return filepath.SkipDir
	}

	if filepath.Ext(path) == ".go" {
		if info.ModTime().After(w.b.LastBuild()) {
			p, _ := filepath.Rel(w.dir, path)
			b_logger.Logger().Info("[", w.b.AppName(), "]", p, "modified.")
			okay := w.b.Build()
			if okay {
				re := w.r.Run()
				return re
			}
		}
	}
	return nil
}

func (w *Watcher) Watch() error {
	b_logger.Logger().Info("Started watching", w.dir)

	// Do a first build
	okay := w.b.Build()
	if okay {
		// Do a first run.
		if re := w.r.Run(); re != nil {
			b_logger.Logger().Error(re.Error())
		}
	}
	stopWatch := make(chan error)
	go func() {
		for {
			err := filepath.Walk(w.dir, w.watchFunc)
			if err != nil && err != filepath.SkipDir {
				stopWatch <- err
				break
			}
			time.Sleep(SLEEP_TIME * time.Millisecond)
		}
	}()
	return <-stopWatch
}
