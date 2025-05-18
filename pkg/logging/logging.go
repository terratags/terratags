package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// ValidLogLevels contains all valid log level options
var ValidLogLevels = []string{"DEBUG", "INFO", "WARN", "ERROR"}

// Initialize sets up the logger with the specified log level
func Initialize(logLevel string) error {
	// Set log format to JSON for better structured logging
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout
	log.SetOutput(os.Stdout)

	// Validate and set log level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: %s. Valid options are: %s", logLevel, strings.Join(ValidLogLevels, ", "))
	}

	return nil
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info logs an info message (only shown in verbose mode)
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Print logs a message that should always be shown regardless of verbosity
func Print(format string, args ...interface{}) {
	// Use Info level but force output
	savedLevel := log.GetLevel()
	log.SetLevel(logrus.InfoLevel)
	log.Infof(format, args...)
	log.SetLevel(savedLevel)
}
