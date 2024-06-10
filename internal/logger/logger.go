// internal/logger/logger.go

package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	return log
}
