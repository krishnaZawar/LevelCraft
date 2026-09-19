//go:build !windows

package entity

import (
	"errors"
	"syscall"
)

// kills the whole process group led by p.Cmd.Process (see executor.setProcessGroup)
//
// ESRCH and EPERM both mean the group is already gone: the group is led by our
// own child, so while that child lives the group exists and is ours to signal.
func killProcessTree(p *Process) error {
	err := syscall.Kill(-p.Cmd.Process.Pid, syscall.SIGKILL)
	if err == nil || errors.Is(err, syscall.ESRCH) || errors.Is(err, syscall.EPERM) {
		return nil
	}
	return err
}
