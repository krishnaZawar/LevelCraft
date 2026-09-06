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

	if err := killProcessTree(p); err != nil {
		if !errors.Is(err, os.ErrProcessDone) {
			return err
		}
	}

	_ = p.Cmd.Wait()
	return nil
}

type CommandConfig struct {
	Pwd  string   // present working directory of the command
	Name string   // name of the executable
	Args []string // args passed to the executable
	Port string   // port on which the process should listen for health check
	Env  []string // extra env vars, added on top of the parent's environment
}
