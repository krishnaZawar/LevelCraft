package component

import (
	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/basecomp"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
)

type ComponentRegistry = models.ComponentRegistry

// contains the registry object registered with all the components base copies
func NewComponentRegistry() *models.ComponentRegistry {
	compRegistry := models.NewComponentRegistry()
	compRegistry.Register(base.ComponentName_Transform, basecomp.NewBaseTransform())
	compRegistry.Register(base.ComponentName_Color, basecomp.NewBaseColor())

	return compRegistry
}
