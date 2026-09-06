//go:build windows

package executor

import "os/exec"

// no-op on Windows; tree termination is handled via taskkill /t in entity.Process.Stop
func setProcessGroup(_ *exec.Cmd) {}
