package app

import (
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

type builderEvents int

const (
	// BUILD is the event to start a build
	BUILD builderEvents = 0
	// TERMINATEBUILDER is the event to TERMINATEBUILDER the build watcher
	TERMINATEBUILDER builderEvents = 1
)

// BuilderOptions provides additional options of the builder
type BuilderOptions struct {
	OnSuccess func([]byte)
	OnFailure func([]byte, error)
}

// Builder builds the project
type Builder struct {
	sync.RWMutex
	lastBuildTime time.Time
	eventChan     chan builderEvents
	Options       BuilderOptions
	app           string
}

// NewBuilder returns a new instance of builder
func NewBuilder(app, dir string, options BuilderOptions) *Builder {
	if runtime.GOOS == "windows" && !strings.HasSuffix(app, ".exe") {
		app += ".exe"
	}

	b := &Builder{
		lastBuildTime: time.Now(),
		eventChan:     make(chan builderEvents),
		Options:       options,
		app:           app,
	}
	go b.Build(dir)
	return b
}

// LastBuildTime retunrs the last time the build process was successful
func (b *Builder) LastBuildTime() time.Time {
	b.RLock()
	defer b.RUnlock()
	return b.lastBuildTime
}

// InitiateBuild initiates a build
func (b *Builder) InitiateBuild() {
	b.eventChan <- BUILD
}

// TerminateBuild terminates a build loop
func (b *Builder) TerminateBuild() {
	b.eventChan <- TERMINATEBUILDER
}

// Build preforms the build
func (b *Builder) Build(dir string) {
	for {
		e := <-b.eventChan
		switch e {
		case BUILD:
			cmd := exec.Command("go", "build", "-o", b.app)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				jww.ERROR.Printf("%s. Cause: %v\n", out, err)
				if b.Options.OnFailure != nil {
					b.Options.OnFailure(out, err)
				}
			} else {
				b.Lock()
				b.lastBuildTime = time.Now()
				b.Unlock()
				if b.Options.OnSuccess != nil {
					b.Options.OnSuccess(out)
				}
			}
		case TERMINATEBUILDER:
			return
		}
	}
}
