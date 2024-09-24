package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
	logDir string
}

func NewLogrusLogger(logDir string) Logger {
	l := logrus.New()
	l.SetOutput(l.Writer())
	l.SetFormatter(&logrus.JSONFormatter{})
	return &LogrusLogger{logger: l, logDir: logDir}
}

// It sets the output file to the current date file
func (l *LogrusLogger) setOutputFileToCurrentDateFile() {
	// check if the log directory exists, if not create it
	if err := os.MkdirAll(l.logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(l.logDir, currentDate+".log")

	// open or create the log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	l.logger.SetOutput(logFile)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.setOutputFileToCurrentDateFile()
	l.logger.Info(args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.setOutputFileToCurrentDateFile()
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.setOutputFileToCurrentDateFile()
	l.logger.Error(args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.setOutputFileToCurrentDateFile()
	l.logger.Debug(args...)
}

func (l *LogrusLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	l.setOutputFileToCurrentDateFile()
	return l.logger.WithFields(fields)
}
