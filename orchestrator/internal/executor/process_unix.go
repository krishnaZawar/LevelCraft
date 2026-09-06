//go:build !windows

package executor

import (
	"os/exec"
	"syscall"
)

// makes the process its own process-group leader, so its whole tree can be killed together
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
