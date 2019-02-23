package logger

import (
	"github.com/sirupsen/logrus"
)

const DEFAULT_LOG_LEVEL = logrus.InfoLevel

var Logger logrus.FieldLogger = logrus.StandardLogger()

func SetLevel(level string) {
	loggerLevel, err := logrus.ParseLevel(level)

	if err != nil {
		Logger.Fatalln("Could not parse log-level: ", err)
	}

	logrus.SetLevel(loggerLevel)
}
