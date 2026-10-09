package input

import (
	_ "embed"
	"fmt"

	"encoding/json"
)

type InputMapping struct {
	Keyboard []InputData `json:"keyboard"`
	Mouse    []InputData `json:"mouse"`
}

// Checks whether the keyName exists in the InputMapping
func (im *InputMapping) KeyExists(keyName string) error {
	for _, key := range im.Keyboard {
		if key.Name == keyName {
			return nil
		}
	}
	for _, key := range im.Mouse {
		if key.Name == keyName {
			return nil
		}
	}
	return fmt.Errorf("\"%s\" key does not exist in the mapping", keyName)
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
