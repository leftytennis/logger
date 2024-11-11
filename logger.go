// Package logger provides a custom log writer that adds a timestamp to each log entry.
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// LogLevel is a type for log levels
type LogLevel int

const (
	// LogLevelNone is the lowest log level and will not log anything
	LogLevelNone LogLevel = iota
	// LogLevelInfo is the log level for info messages
	LogLevelInfo = 1 << iota
	// LogLevelWarn is the log level for warning messages
	LogLevelWarn
	// LogLevelError is the log level for error messages
	LogLevelError
	// LogLevelDebug is the log level for debug messages
	LogLevelDebug
	// LogLevelTrace is the log level for trace messages
	LogLevelTrace
)

const (
	// LogDateFormat is the format for the timestamp of log entries
	LogDateFormat string = "2006-01-02 15:04:05.000 MST"
)

// Logger is a custom log writer that adds a timestamp to each log entry
type Logger struct {
	Level LogLevel
	File  *os.File
	m     sync.RWMutex
}

// Options are options for the Logger
type Options struct {
	Level LogLevel
	File  *os.File
}

func (l LogLevel) String() string {
	switch l {
	case LogLevelNone:
		return "none"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	case LogLevelDebug:
		return "debug"
	case LogLevelTrace:
		return "trace"
	default:
		return "unknown"
	}
}

// New creates a new Logger
func New() *Logger {
	return &Logger{Level: LogLevelInfo, File: os.Stderr}
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(opts Options) *Logger {
	if opts.Level == 0 {
		opts.Level = LogLevelInfo
	}
	if opts.File == nil {
		opts.File = os.Stderr
	}
	return &Logger{Level: opts.Level, File: opts.File}
}

// SetLevel sets the log level
func (writer *Logger) SetLevel(level LogLevel) {
	writer.Level = level
	return
}

// SetOutput sets the output file for the logger
func (writer *Logger) SetOutput(file *os.File) {
	if file == nil {
		writer.File = os.Stderr
	} else {
		writer.File = file
	}
	return
}

// Write writes a log entry to an output file (default: os.Stderr)
func (writer *Logger) Write(bytes []byte) (int, error) {

	if writer.File == nil {
		panic("file is nil")
	}

	writer.m.Lock()
	defer writer.m.Unlock()

	if bytes[len(bytes)-1] != '\n' {
		bytes = append(bytes, '\n')
	}

	return writer.File.Write(bytes)
}

// Debug logs a debug message
func (writer *Logger) Debug(message string) {
	if writer.Level >= LogLevelDebug {
		message := fmt.Sprintf(time.Now().Format(LogDateFormat) + " D " + message + "\n")
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Debugln logs a debug message with a newline
func (writer *Logger) Debugln(message string) {
	writer.Debug(message)
	return
}

// Debugf logs a debug message with a format string
func (writer *Logger) Debugf(format string, args ...interface{}) {
	if writer.Level >= LogLevelDebug {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" D "+format, args...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Error logs an error message
func (writer *Logger) Error(message string) {
	if writer.Level >= LogLevelError {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat) + " E " + message + "\n")
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Errorln logs an error message with a newline
func (writer *Logger) Errorln(message string) {
	writer.Error(message)
	return
}

// Errorf logs an error message with a format string
func (writer *Logger) Errorf(format string, args ...interface{}) {
	if writer.Level >= LogLevelError {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" E "+format, args...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Fatal logs a fatal message
func (writer *Logger) Fatal(message string) {
	msg := fmt.Sprintf(time.Now().Format(LogDateFormat) + " F " + message + "\n")
	_, err := writer.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

// Fatalln logs a fatal message with a newline
func (writer *Logger) Fatalln(message string) {
	writer.Fatal(message)
	return
}

// Fatalf logs a fatal message with a format string
func (writer *Logger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" F "+format, args...)
	_, err := writer.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

// Info logs an info message
func (writer *Logger) Info(message string) {
	if writer.Level >= LogLevelInfo {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat) + " I " + message + "\n")
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Infoln logs an info message with a newline
func (writer *Logger) Infoln(message string) {
	writer.Info(message)
	return
}

// Infof logs an info message with a format string
func (writer *Logger) Infof(format string, args ...interface{}) {
	if writer.Level >= LogLevelInfo {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" I "+format, args...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Warn logs a warning message
func (writer *Logger) Warn(message string) {
	if writer.Level >= LogLevelWarn {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat) + " W " + message + "\n")
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Warnln logs a warning message with a newline
func (writer *Logger) Warnln(message string) {
	writer.Warn(message)
	return
}

// Warnf logs a warning message with a format string
func (writer *Logger) Warnf(format string, args ...interface{}) {
	if writer.Level >= LogLevelWarn {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" W "+format, args...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}
