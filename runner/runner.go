package runner

import (
	"os/exec"
	"time"
	"path/filepath"
	"os"
)

/**
 * Runner struct.
 * proc: absolute path of process to run.
 * command: run command
 */
type Runner struct {
	proc string
	command *exec.Cmd
}

/**
 * Returns a new runner wi
 */
func NewRunner(name, dir string) Runner {
	r := Runner{proc: filepath.Join(dir, name)}
	return r
}

func (r *Runner) Run() error {
	er := r.terminateIfRunning()
	if er != nil {
		return er
	}
	r.command = exec.Command(r.proc)
	r.command.Stderr = os.Stderr
	r.command.Stdout = os.Stdout
	return r.command.Start()
}

func (r *Runner) terminateIfRunning() (error) {
	if r.command == nil {
		return nil
	}
	timer := time.NewTimer(time.Second)
	done := make(chan error, 1)
	go func() {
		done <- r.command.Wait()
	}()

	select {
	case <-timer.C:
		r.command.Process.Kill()
		<-done
	case err := <-done:
		if err != nil {
			return err
		}
	}
	r.command = nil
	return nil
}
