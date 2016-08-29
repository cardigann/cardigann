package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	Logger logrus.FieldLogger
)

func init() {
	Logger = logrus.New()
	Logger.(*logrus.Logger).Level = logrus.DebugLevel
	Logger.(*logrus.Logger).Out = os.Stderr
}

func SetLevel(level logrus.Level) {
	Logger.(*logrus.Logger).Level = level
}
