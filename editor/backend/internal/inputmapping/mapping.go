package inputmapping

import "github.com/krishnaZawar/LevelCraft/utils/input"

var mapping *input.InputMapping

func Load() error {
	var err error
	mapping, err = input.LoadMapping()
	if err != nil {
		return err
	}
	return nil
}

func Get() *input.InputMapping {
	return mapping
}
