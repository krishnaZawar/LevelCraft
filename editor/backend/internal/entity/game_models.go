package entity

type SaveGameRequest struct {
	Filepath string `json:"filepath"` // destination where the game state should be saved. (JSON file)
}

type SaveGameResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetGameStateResponse struct {
	Success   bool                   `json:"success"`
	GameState map[string]interface{} `json:"gameState"`
}

type LoadGameRequest struct {
	Filepath string `json:"filepath"` // source to load the game state from. (JSON file)
}

type LoadGameResponse struct {
	Success   bool                   `json:"success"`
	GameState map[string]interface{} `json:"gameState"`
}
