//go:build !windows

package entity

import (
	"errors"
	"syscall"
)

// kills the whole process group led by p.Cmd.Process (see executor.setProcessGroup)
func killProcessTree(p *Process) error {
	pgid := p.Cmd.Process.Pid
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	return nil
}
