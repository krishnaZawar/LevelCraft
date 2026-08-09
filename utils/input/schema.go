package input

import (
	_ "embed"

	"encoding/json"
)

type InputMapping struct {
	Keyboard []InputData `json:"keyboard"`
	Mouse    []InputData `json:"mouse"`
}

type InputData struct {
	Name  string `json:"name"`
	Code  int    `json:"code"`
	Label string `json:"label"`
}

//go:embed mapping.json
var inputJSON []byte

func LoadMapping() (*InputMapping, error) {
	var mapping *InputMapping
	if err := json.Unmarshal(inputJSON, &mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
