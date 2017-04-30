package app

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

// ChangableWriter is a custom writer with a changable destination
type ChangableWriter struct {
	sync.Mutex
	io.Writer
}

// ChangableReader is a custom reader with a changable source
type ChangableReader struct {
	sync.Mutex
	io.Reader
}

func (c *ChangableWriter) Write(p []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	return c.Writer.Write(p)
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
		Err:     &ChangableWriter{Writer: os.Stderr},
		Out:     &ChangableWriter{Writer: os.Stdout},
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
func (r *Runner) SetErr(w io.Writer) {
	r.Err.Lock()
	defer r.Err.Unlock()
	r.Err.Writer = w
}

// SetOut sets the std out of the command
func (r *Runner) SetOut(w io.Writer) {
	r.Out.Lock()
	defer r.Out.Unlock()
	r.Out.Writer = w
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
