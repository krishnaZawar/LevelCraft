package entity

type UpdateComponentRequest struct {
	Details map[string]interface{} `json:"details"` // new details of the component
}

type ComponentResponse struct {
	Success       bool                   `json:"success"`
	ObjectDetails map[string]interface{} `json:"objectDetails"`
}
