package app

import (
	"io"
	"os"
	"os/exec"
	"sync"

	jww "github.com/spf13/jwalterweatherman"
)

// ChangableWriter is a custom writer with a changable destination
type ChangableWriter struct {
	sync.Mutex
	io.Writer
	ws map[string]io.Writer
}

// ChangableReader is a custom reader with a changable source
type ChangableReader struct {
	sync.Mutex
	io.Reader
}

func (c *ChangableWriter) Write(p []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	b, err := c.Writer.Write(p)
	if err != nil {
		return b, err
	}
	for id, w := range c.ws {
		_, err := w.Write(p)
		if err != nil {
			jww.ERROR.Println(err)
			delete(c.ws, id)
		}
	}
	return b, nil
}

// Runner defines the a runner
type Runner struct {
	cmd     *exec.Cmd
	command string
	args    []string
	Err     *ChangableWriter
	Out     *ChangableWriter
	In      *ChangableReader
}

// NewRunner returns a new runner
func NewRunner(command string, args ...string) *Runner {
	return &Runner{
		Err:     &ChangableWriter{Writer: os.Stderr, ws: make(map[string]io.Writer)},
		Out:     &ChangableWriter{Writer: os.Stdout, ws: make(map[string]io.Writer)},
		In:      &ChangableReader{Reader: os.Stdin},
		command: command,
		args:    args,
	}
}

// Run runs the build binary
func (r *Runner) Run() error {
	if err := r.Terminate(); err != nil {
		return err
	}
	cmd := exec.Command(r.command, r.args...)
	cmd.Stderr = r.Err
	cmd.Stdout = r.Out
	r.cmd = cmd
	return r.cmd.Start()
}

// SetErr sets the std err of the command
func (r *Runner) SetErr(w io.Writer, id string) {
	r.Err.Lock()
	defer r.Err.Unlock()
	r.Err.ws[id] = w
}

// SetOut sets the std out of the command
func (r *Runner) SetOut(w io.Writer, id string) {
	r.Out.Lock()
	defer r.Out.Unlock()
	r.Out.ws[id] = w
}

// UnsetErr unsets the std err of the command
func (r *Runner) UnsetErr(id string) {
	r.Err.Lock()
	defer r.Err.Unlock()
	delete(r.Err.ws, id)
}

// UnsetOut unsets the std out of the command
func (r *Runner) UnsetOut(id string) {
	r.Out.Lock()
	defer r.Out.Unlock()
	delete(r.Out.ws, id)
}

// SetIn sets the std in of the command
func (r *Runner) SetIn(rr io.Reader) {
	r.In.Lock()
	defer r.In.Unlock()
	r.In.Reader = rr
}

// Terminate terminates the running process
func (r *Runner) Terminate() error {
	if r.cmd == nil {
		return nil
	}
	if err := r.cmd.Process.Kill(); err != nil {
		return err
	}
	if err := r.cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	return nil
}
