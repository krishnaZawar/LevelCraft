package executor

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
	"github.com/krishnaZawar/LevelCraft/utils/logger"
)

var ls = logger.New(base.ServiceName)
var client = &http.Client{
	Timeout: base.PingTimeout,
}

// generates a random port number for the process to bind to
func getRandomPort() string {
	return strconv.Itoa(rand.IntN(base.MaxPortValue+1-base.MinPortValue) + base.MinPortValue)
}

// pings the process endpoint for health check
// fails if the statusCode was not 200 or error was encountered
func CheckProcessHealth(baseUrl string) error {
	resp, err := client.Get(baseUrl + "/ping")
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code not found. expected status code: %d recieved status code: %d", http.StatusOK, resp.StatusCode)
	}
	return nil
}

// builds and runs the comm command and returns the Process details and error if any
func buildAndRunProcess(processName string, comm *entity.CommandConfig) (*entity.Process, error) {
	cmd := exec.Command(comm.Name, comm.Args...)
	cmd.Dir = comm.Pwd
	cmd.Stdout = os.Stdout
	if len(comm.Env) > 0 {
		cmd.Env = append(os.Environ(), comm.Env...)
	}
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &entity.Process{
		Name:             processName,
		Cmd:              cmd,
		CommunicationURI: "http://localhost:" + comm.Port,
	}, nil
}

// polls /ping until it succeeds or timeout elapses
func waitForHealthy(baseUrl string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := CheckProcessHealth(baseUrl); err == nil {
			return true
		}
		time.Sleep(base.PingPollInterval)
	}
	return false
}

// spawns a process via buildComm and waits for it to become healthy, retrying up to base.MaxStartAttempts times
func startAndWaitHealthy(
	processName string,
	buildComm func() entity.CommandConfig,
	startupTimeout time.Duration,
) (*entity.Process, error) {
	var lastErr error
	for attempt := 1; attempt <= base.MaxStartAttempts; attempt++ {
		ls.Info().Msgf("starting with attempt %d of starting %s", attempt, processName)

		comm := buildComm()
		process, err := buildAndRunProcess(processName, &comm)
		if err != nil {
			lastErr = err
			ls.ErrorWith(err).Msgf("failed to start %s", processName)
			continue
		}

		if waitForHealthy(process.CommunicationURI, startupTimeout) {
			ls.Info().Msgf("%s started... PID: %d", processName, process.Cmd.Process.Pid)
			return process, nil
		}

		lastErr = fmt.Errorf("%s did not become healthy within %s", processName, startupTimeout)
		ls.Warn().Msgf("%s failed to become healthy within %s, stopping and retrying", processName, startupTimeout)
		if stopErr := process.Stop(); stopErr != nil {
			ls.ErrorWith(stopErr).Msgf("failed to stop unhealthy %s process: %d", processName, process.Cmd.Process.Pid)
		}
	}

	return nil, fmt.Errorf("failed to start %s after %d attempts: %w", processName, base.MaxStartAttempts, lastErr)
}
