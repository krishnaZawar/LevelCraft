package entity

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func startSleeper(t *testing.T) *Process {
	t.Helper()
	cmd := exec.Command("sleep", "30")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start test process: %v", err)
	}
	return &Process{Name: "test", Cmd: cmd}
}

// A zombie is a process that has exited but was never waited on, which is what
// "still there" means for a pid nothing is running under any more.
func processState(t *testing.T, pid int) string {
	t.Helper()
	out, err := exec.Command("ps", "-o", "stat=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return "" // the pid is gone entirely, which is the outcome we want
	}
	return strings.TrimSpace(string(out))
}

func isZombie(t *testing.T, pid int) bool {
	t.Helper()
	return strings.HasPrefix(processState(t, pid), "Z")
}

func Test_Stop_RunningProcess(t *testing.T) {
	p := startSleeper(t)
	pid := p.Cmd.Process.Pid

	if err := p.Stop(); err != nil {
		t.Fatalf("Stop() returned error: %v", err)
	}

	if isZombie(t, pid) {
		t.Fatalf("Stop() left the process as a zombie (state %q)", processState(t, pid))
	}
}

// The path that matters for a game run: the player window is closed by the
// user, so the process is already gone by the time the orchestrator stops it.
// Killing an empty process group can report EPERM rather than ESRCH, and an
// early return there would skip the wait and leak a zombie.
func Test_Stop_AlreadyExitedProcess(t *testing.T) {
	p := startSleeper(t)
	pid := p.Cmd.Process.Pid

	// Kill it out from under the orchestrator, exactly as closing the window would.
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
		t.Fatalf("failed to kill the test process: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	if err := p.Stop(); err != nil {
		t.Fatalf("Stop() returned error for an already exited process: %v", err)
	}

	if isZombie(t, pid) {
		t.Fatalf("Stop() left an already exited process as a zombie (state %q)", processState(t, pid))
	}
}

func Test_Stop_NeverStarted(t *testing.T) {
	if err := (&Process{Name: "test"}).Stop(); err != nil {
		t.Fatalf("Stop() on an unstarted process returned error: %v", err)
	}
}
