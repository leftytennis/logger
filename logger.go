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
	// LogLevelFatal is the highest log level and will log fatal messages
	LogLevelFatal
)

const (
	// LogDateFormat is the format for the timestamp of log entries
	LogDateFormat string = "2006-01-02 15:04:05.000 MST"
)

// Logger is a custom log writer that adds a timestamp to each log entry
type Logger struct {
	Level  LogLevel
	Output *os.File
	m      sync.RWMutex
}

// Options are options for the Logger
type Options struct {
	Level  LogLevel
	Output *os.File
}

func (l LogLevel) String() string {
	switch l {
	case LogLevelNone:
		return " "
	case LogLevelFatal:
		return "F"
	case LogLevelInfo:
		return "I"
	case LogLevelWarn:
		return "W"
	case LogLevelError:
		return "E"
	case LogLevelDebug:
		return "D"
	case LogLevelTrace:
		return "T"
	default:
		return "U"
	}
}

func buildString(l LogLevel, a ...any) string {
	message := fmt.Sprintf(time.Now().Format(LogDateFormat)) + " " + l.String()
	for _, v := range a {
		message += " " + v.(string)
	}
	return message
}

// New creates a new Logger
func New() *Logger {
	return &Logger{Level: LogLevelInfo, Output: os.Stderr}
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(opts Options) *Logger {
	if opts.Level == 0 {
		opts.Level = LogLevelInfo
	}
	if opts.Output == nil {
		opts.Output = os.Stderr
	}
	return &Logger{Level: opts.Level, Output: opts.Output}
}

// SetLevel sets the log level
func (writer *Logger) SetLevel(level LogLevel) {
	writer.Level = level
	return
}

// SetOutput sets the output file for the logger
func (writer *Logger) SetOutput(file *os.File) {
	if file == nil {
		writer.Output = os.Stderr
	} else {
		writer.Output = file
	}
	return
}

// Write writes a log entry to an output file (default: os.Stderr)
func (writer *Logger) Write(bytes []byte) (int, error) {

	if writer.Output == nil {
		panic("file is nil")
	}

	writer.m.Lock()
	defer writer.m.Unlock()

	if bytes[len(bytes)-1] != '\n' {
		bytes = append(bytes, '\n')
	}

	return writer.Output.Write(bytes)
}

// Debug logs a debug message
func (writer *Logger) Debug(a ...any) {
	if writer.Level >= LogLevelDebug {
		message := buildString(LogLevelDebug, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Debugln logs a debug message with a newline
func (writer *Logger) Debugln(a ...any) {
	writer.Debug(a...)
	return
}

// Debugf logs a debug message with a format string
func (writer *Logger) Debugf(format string, a ...any) {
	if writer.Level >= LogLevelDebug {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" D "+format, a...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Error logs an error message
func (writer *Logger) Error(a ...any) {
	if writer.Level >= LogLevelError {
		message := buildString(LogLevelError, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Errorln logs an error message with a newline
func (writer *Logger) Errorln(a ...any) {
	writer.Error(a...)
	return
}

// Errorf logs an error message with a format string
func (writer *Logger) Errorf(format string, a ...any) {
	if writer.Level >= LogLevelError {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" E "+format, a...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Fatal logs a fatal message
func (writer *Logger) Fatal(a ...any) {
	message := buildString(LogLevelFatal, a...)
	_, err := writer.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

// Fatalln logs a fatal message with a newline
func (writer *Logger) Fatalln(a ...any) {
	writer.Fatal(a...)
	return
}

// Fatalf logs a fatal message with a format string
func (writer *Logger) Fatalf(format string, a ...any) {
	msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" F "+format, a...)
	_, err := writer.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

// Info logs an info message
func (writer *Logger) Info(a ...any) {
	if writer.Level >= LogLevelInfo {
		message := buildString(LogLevelInfo, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Infoln logs an info message with a newline
func (writer *Logger) Infoln(a ...any) {
	writer.Info(a...)
	return
}

// Infof logs an info message with a format string
func (writer *Logger) Infof(format string, a ...any) {
	if writer.Level >= LogLevelInfo {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" I "+format, a...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Warn logs a warning message
func (writer *Logger) Warn(a ...any) {
	if writer.Level >= LogLevelWarn {
		message := buildString(LogLevelWarn, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
	return
}

// Warnln logs a warning message with a newline
func (writer *Logger) Warnln(a ...any) {
	writer.Warn(a...)
	return
}

// Warnf logs a warning message with a format string
func (writer *Logger) Warnf(format string, a ...any) {
	if writer.Level >= LogLevelWarn {
		msg := fmt.Sprintf(time.Now().Format(LogDateFormat)+" W "+format, a...)
		_, err := writer.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}
	return
}
