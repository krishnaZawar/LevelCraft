package entity

type CreateGameobjectResponse struct {
	Success       bool                   `json:"success"`
	ObjectDetails map[string]interface{} `json:"objectDetails"`
}

// UpdateGameobjectRequest carries the gameobject metadata to change.
// Fields are pointers so an omitted one is left alone rather than blanked.
type UpdateGameobjectRequest struct {
	Name  *string `json:"name"`
	Group *string `json:"group"`
}

type UpdateGameobjectResponse struct {
	Success       bool                   `json:"success"`
	ObjectDetails map[string]interface{} `json:"objectDetails"`
}

type DuplicateGameobjectResponse struct {
	Success       bool                   `json:"success"`
	ObjectDetails map[string]interface{} `json:"objectDetails"`
}

type DeleteGameobjectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
