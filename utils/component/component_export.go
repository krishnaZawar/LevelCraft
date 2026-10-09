package component

import (
	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/basecomp"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/movement"
	"github.com/krishnaZawar/LevelCraft/utils/input"
)

// holds the list of all the components that will be used by the editor to display all the components available
var (
	ComponentList = []string{
		base.ComponentName_Transform,
		base.ComponentName_Color,
		base.ComponentName_InputMovement,
	}
)

type Component = models.Component
type ComponentDetails = models.ComponentDetails

type Transform = basecomp.Transform

func NewTransform(x, y, w, h int) *Transform {
	return basecomp.NewTransform(x, y, w, h)
}

type Color = basecomp.Color

func NewColor(r, g, b, a int) *Color {
	return basecomp.NewColor(r, g, b, a)
}

type InputMovement = movement.InputMovement

func NewInputMovement(inputmapping *input.InputMapping, speed int) *movement.InputMovement {
	return movement.NewInputMovement(inputmapping, speed)
}
