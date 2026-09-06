//go:build windows

package entity

import (
	"os/exec"
	"strconv"
)

// kills the process tree by PID (no process-group concept on Windows)
func killProcessTree(p *Process) error {
	cmd := exec.Command("taskkill", "/pid", strconv.Itoa(p.Cmd.Process.Pid), "/t", "/f")
	return cmd.Run()
}
