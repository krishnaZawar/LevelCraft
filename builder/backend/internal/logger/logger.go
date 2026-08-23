package logger

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/utils/logger"
)

var ls = logger.New(base.ServiceName)

func Get() *logger.Logger {
	return ls
}
