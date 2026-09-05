package executor

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"strconv"

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

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &entity.Process{
		Name:             processName,
		Cmd:              cmd,
		CommunicationURI: "http://localhost:" + comm.Port,
	}, nil
}
