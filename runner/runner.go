package runner

import (
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

/**
 * Runner struct.
 * proc    : absolute path of process to run.
 * command : run command
 */
type Runner struct {
    proc    string
    command *exec.Cmd
}

/**
 * Returns a new runner.
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
            return err
        }
    }
    r.command = nil
    return nil
}
