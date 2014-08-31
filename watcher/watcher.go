package watcher

import (
	"github.com/v4run/bob/bLogger"
	"github.com/v4run/bob/builder"
	"github.com/v4run/bob/runner"
	"os"
	"path/filepath"
	"time"
	"strings"
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
	b   builder.Builder
	r   runner.Runner
	dir string
}

/**
 * Returns a new watcher.
 */
func NewWatcher(dir string) Watcher {
	appName := filepath.Base(dir)
	return Watcher{dir: dir, b: builder.NewBuilder(appName, dir), r: runner.NewRunner(appName, dir)}
}

func (w *Watcher) watchFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if strings.HasPrefix(path, ".") {  // skip directories like .git, .idea etc.
		return filepath.SkipDir
	}

	if filepath.Ext(path) == ".go" {
		if info.ModTime().After(w.b.LastBuild()) {
			bLogger.Logger().Info(path, "modified. Rebuilding.")
			out, er := w.b.Build()
			bLogger.Logger().Info("Build complete.")
			if er != nil {
				bLogger.Logger().Error(string(out))
			} else {
				re := w.r.Run()
				return re
			}
		}
	}
	return nil
}

func (w *Watcher) Watch() error {
	out, er := w.b.Build() // Do a first build.
	if er != nil {
		bLogger.Logger().Error(string(out))
		return er
	}
	bLogger.Logger().Info("Started watching", w.dir)
	stopWatch := make(chan error)
	go func() {
		for {
			err := filepath.Walk(w.dir, w.watchFunc)
			if err != nil {
				stopWatch <- err
			}
			time.Sleep(SLEEP_TIME * time.Millisecond)
		}
	}()
	return <-stopWatch
}
