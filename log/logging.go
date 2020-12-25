package log

import (
	"github.com/reed/blockchain/config"
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func Init() {
	fileHooker := NewFileRotateHooker(config.LogDir(), config.Default.LogAge)

	Logger = logrus.New()
	Logger.Hooks.Add(fileHooker)
	Logger.Out = os.Stdout
	Logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	Logger.Level = logrus.DebugLevel
}
