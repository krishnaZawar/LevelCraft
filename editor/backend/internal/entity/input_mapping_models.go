package entity

import (
	"github.com/krishnaZawar/LevelCraft/utils/input"
)

type GetInputMappingResponse struct {
	Success bool               `json:"success"`
	Mapping input.InputMapping `json:"mapping"`
}
