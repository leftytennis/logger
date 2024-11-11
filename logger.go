// Package logger provides a custom log writer that adds a timestamp to each log entry.
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	LogLevelNone  LogLevel = iota
	LogLevelFatal          = 1 << iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelDebug
	LogLevelTrace
)

const (
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
	case LogLevelFatal:
		return "fatal"
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
}

// SetOutput sets the output file for the logger
func (writer *Logger) SetOutput(file *os.File) {
	if file == nil {
		writer.File = os.Stderr
	} else {
		writer.File = file
	}
}

// Write writes a log entry to stderr
func (writer *Logger) Write(bytes []byte) (int, error) {

	writer.m.Lock()
	defer writer.m.Unlock()

	if bytes[len(bytes)-1] != '\n' {
		bytes = append(bytes, '\n')
	}

	// tempStr := string(bytes)

	// if !strings.HasSuffix(tempStr, "\n") {
	// 	tempStr += "\n"
	// }

	if writer.File == nil {
		panic("file is nil")
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
}

// Debugln logs a debug message with a newline
func (writer *Logger) Debugln(message string) {
	writer.Debug(message)
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
}

// Errorln logs an error message with a newline
func (writer *Logger) Errorln(message string) {
	writer.Error(message)
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
}

// Infoln logs an info message with a newline
func (writer *Logger) Infoln(message string) {
	writer.Info(message)
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
}

// Warnln logs a warning message with a newline
func (writer *Logger) Warnln(message string) {
	writer.Warn(message)
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
}
