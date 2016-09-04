package logger

import (
	"io"
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

func SetFormatter(f logrus.Formatter) {
	Logger.(*logrus.Logger).Formatter = f
}

func SetOutput(out io.Writer) {
	Logger.(*logrus.Logger).Out = out
}

func SetLevel(level logrus.Level) {
	Logger.(*logrus.Logger).Level = level
}

func AddHook(h logrus.Hook) {
	Logger.(*logrus.Logger).Hooks.Add(h)
}
