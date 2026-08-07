package entity

type CreateGameobjectResponse struct {
	Success       bool                   `json:"success"`
	ObjectDetails map[string]interface{} `json:"objectDetails"`
}

type DeleteGameobjectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
