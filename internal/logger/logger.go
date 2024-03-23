package logger

import "github.com/sirupsen/logrus"

func InitLogger(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Warnf("Can't parse %s as level, using default Info", level)

		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)
}
