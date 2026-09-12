package entity

import (
	"errors"
	"os"
	"os/exec"
)

type Process struct {
	Name             string    // name of the process being instantiated
	Cmd              *exec.Cmd // command that is being executed
	CommunicationURI string    // communication url of the process that will be passed to others on startup for communication
}

// forcefully terminates the process and its whole process tree (see killProcessTree)
func (p *Process) Stop() error {
	if p.Cmd == nil || p.Cmd.Process == nil {
		return nil
	}

	killErr := killProcessTree(p)

	// Waited for unconditionally, even when the kill reported a failure: a
	// process that already exited on its own stays a zombie until it is reaped.
	_ = p.Cmd.Wait()

	if killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
		return killErr
	}
	return nil
}

type CommandConfig struct {
	Pwd  string   // present working directory of the command
	Name string   // name of the executable
	Args []string // args passed to the executable
	Port string   // port on which the process should listen for health check
	Env  []string // extra env vars, added on top of the parent's environment
}
