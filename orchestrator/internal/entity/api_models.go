package entity

// RunGameRequest is the payload the editor sends to start a game run
type RunGameRequest struct {
	Filepath string `json:"filepath"` // the scene file the builder should run (JSON file)
}

// GameActionResponse is returned for the run/stop lifecycle requests
type GameActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GameStatusResponse tells the editor whether a game is currently running.
// Polled, so a run that ended on its own still resets the editor UI
type GameStatusResponse struct {
	Success bool `json:"success"`
	Running bool `json:"running"`
}

func NewGameActionResponse(success bool, message string) GameActionResponse {
	return GameActionResponse{
		Success: success,
		Message: message,
	}
}
