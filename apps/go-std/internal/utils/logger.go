package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger is a structured logger with log levels
type Logger struct {
	level     LogLevel
	logger    *log.Logger
	showDebug bool
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, showDebug bool) *Logger {
	return &Logger{
		level:     level,
		logger:    log.New(os.Stdout, "", 0),
		showDebug: showDebug,
	}
}

// getCallerInfo returns the file and line number of the caller
func (l *Logger) getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // Skip 2 frames to get the actual caller
	if !ok {
		return "unknown:0"
	}
	// Get just the file name without the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

// formatMessage formats the log message with timestamp, level, caller info, and message
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := levelNames[level]
	caller := l.getCallerInfo()
	message := fmt.Sprintf(format, args...)

	return fmt.Sprintf("%s [%s] %s - %s", timestamp, levelStr, caller, message)
}

// log writes a log message if the level is appropriate
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Skip debug messages if debug is disabled
	if level == DEBUG && !l.showDebug {
		return
	}

	message := l.formatMessage(level, format, args...)
	l.logger.Println(message)

	// For FATAL level, exit the program
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// SetLevel changes the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetDebug enables or disables debug logging
func (l *Logger) SetDebug(showDebug bool) {
	l.showDebug = showDebug
}
