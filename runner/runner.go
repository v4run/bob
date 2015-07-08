package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/v4run/bob/blogger"
)

/**
 * Runner struct.
 * proc    : absolute path of process to run.
 * command : run command
 */
type Runner struct {
	proc     string
	command  *exec.Cmd
	commands []string
}

/**
 * Returns a new runner.
 */
func NewRunner(name, dir, commands string) Runner {
	r := Runner{proc: filepath.Join(dir, name), commands: prepareCustomCommands(commands)}
	return r
}

func prepareCustomCommands(c string) []string {
	tc := strings.TrimSpace(c)
	if tc == "" {
		return make([]string, 0)
	}
	return strings.Split(tc, " ")
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

func (r *Runner) RunCustom() error {
	if len(r.commands) == 0 {
		return nil
	}
	blogger.Info().Command("running").Message("'" + strings.Join(r.commands, " ") + "'").Log()
	if len(r.commands) == 1 {
		c := exec.Command(r.commands[0])
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Start()
	}
	c := exec.Command(r.commands[0], r.commands[1:]...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Start()
}

func (r *Runner) terminateIfRunning() error {
	if r.command == nil {
		return nil
	}
	timer := time.NewTimer(time.Second)
	done := make(chan error, 1)

	// Try to stop the process.
	go func() {
		done <- r.command.Wait()
	}()

	// Wait for atmost 1 second for the process to terminate.
	// If the process didn't exit within 1 second, kill it.
	select {
	case <-timer.C:
		r.command.Process.Kill()
		<-done
	case err := <-done:
		if err != nil {
			_, ok := err.(*exec.ExitError)
			if !ok {
				return err
			}
		}
	}
	r.command = nil
	return nil
}
