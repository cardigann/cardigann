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
	Logger.(*logrus.Logger).Level = logrus.InfoLevel
	Logger.(*logrus.Logger).Out = os.Stderr

	if os.Getenv("DEBUG") != "" || os.Getenv("CARDIGANN_DEBUG") != "" {
		Logger.(*logrus.Logger).Level = logrus.DebugLevel
	}
}

func SetLevel(level logrus.Level) {
	Logger.(*logrus.Logger).Level = level
}
