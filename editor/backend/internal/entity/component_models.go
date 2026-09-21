package entity

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
)

type UpdateComponentRequest struct {
	Details json.RawMessage `json:"details"` // new details of the component
}

type ComponentResponse struct {
	Success       bool                         `json:"success"`
	ObjectDetails gameobject.GameobjectDetails `json:"objectDetails"`
}
