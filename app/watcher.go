package app

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

const (
	// SLEEPTIME is the time between successive checks for modification
	SLEEPTIME = 400
)

// WatcherEvent defines the type for events in watcher
type WatcherEvent int

const (
	// STOPWATCHING is triggered to stop the watch process
	STOPWATCHING WatcherEvent = 0
)

// Watcher watches the directory and notifies any change
type Watcher struct {
	Dirs        []string
	R           *Runner
	B           *Builder
	EventChan   chan WatcherEvent
	shouldWatch func(os.FileInfo) (bool, error)
}

/**
 * Called for each file in the directory.
 * If a change is found, a build is performed and the service is run.
 */
func (w *Watcher) watchFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if watch, err := w.shouldWatch(info); !watch {
		return err
	}

	if info.ModTime().After(w.B.LastBuildTime()) {
		w.B.InitiateBuild()
	}
	return nil
}

// Watch watches the directory for changes and initiates a build
func (w *Watcher) Watch() error {
	ticker := time.NewTicker(time.Millisecond * SLEEPTIME)
	for {
		select {
		case e := <-w.EventChan:
			switch e {
			case STOPWATCHING:
				w.B.TerminateBuild()
				w.R.Terminate()
				return nil
			}
		default:
			for _, dir := range w.Dirs {
				err := filepath.Walk(dir, w.watchFunc)
				if err != nil && err != filepath.SkipDir {
					return err
				}
				<-ticker.C
			}
			<-ticker.C
		}
	}
}

// NewWatcher returns a new instance of watcher
func NewWatcher(dirs []string, app string) *Watcher {
	if len(dirs) == 0 {
		return nil
	}
	r := NewRunner(filepath.Join(dirs[0], app))
	bo := BuilderOptions{
		OnSuccess: func(out []byte) {
			if err := r.Run(); err != nil {
				jww.ERROR.Println("Run failed", err)
			}
		},
		OnFailure: func(out []byte, err error) {
			jww.ERROR.Println("Build failed", string(out), err)
		},
	}
	b := NewBuilder(app, dirs[0], bo)
	w := &Watcher{
		Dirs: dirs,
		B:    b,
		R:    r,
		shouldWatch: func(path os.FileInfo) (bool, error) {
			if strings.HasPrefix(filepath.Base(path.Name()), ".") {
				if path.IsDir() {
					return false, filepath.SkipDir
				}
				return false, nil
			}
			if filepath.Ext(path.Name()) == ".go" {
				return true, nil
			}
			return false, nil
		},
		EventChan: make(chan WatcherEvent),
	}
	w.B.InitiateBuild()
	return w
}
