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
		base.ComponentName_Patrol,
		base.ComponentName_Follow,
	}
)

type Component = models.Component
type ComponentDetails = models.ComponentDetails

type Transform = basecomp.Transform

func NewTransform(x, y, w, h float64) *Transform {
	return basecomp.NewTransform(x, y, w, h)
}

type Color = basecomp.Color

func NewColor(r, g, b, a int) *Color {
	return basecomp.NewColor(r, g, b, a)
}

type InputMovement = movement.InputMovement

func NewInputMovement(inputmapping *input.InputMapping, speed float64) *movement.InputMovement {
	return movement.NewInputMovement(inputmapping, speed)
}

type Patrol = movement.Patrol
type Point = movement.Point

func NewPatrol(speed float64, points []Point, pingPong bool) *Patrol {
	return movement.NewPatrol(speed, points, pingPong)
}

type Follow = movement.Follow

const (
	FollowTarget_Object = movement.FollowTarget_Object
	FollowTarget_Point  = movement.FollowTarget_Point
)

func NewFollow(speed float64, targetMode string, stopDistance float64) *Follow {
	return movement.NewFollow(speed, targetMode, stopDistance)
}
