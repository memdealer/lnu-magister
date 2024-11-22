package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Configure logging

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info" // default to info if environment variable is not set
	}

	// Parse string to logrus.Level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Error("Invalid LOG_LEVEL, defaulting to INFO")
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)

}
