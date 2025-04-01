// Package logger provides a custom log writer that adds a timestamp to each log entry.
package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel is a type for log levels
type LogLevel int

const (
	// LogLevelFatal is the highest log level and will log fatal messages
	LogLevelFatal = -16
	// LogLevelError is the log level for error messages
	LogLevelError = -8
	// LogLevelWarn is the log level for warning messages
	LogLevelWarn = -4
	// LogLevelInfo is the log level for info messages
	LogLevelInfo = 0
	// LogLevelVerbose is the log level for verbose messages
	LogLevelVerbose = 4
	// LogLevelDebug is the log level for debug messages
	LogLevelDebug = 8
	// LogLevelTrace is the log level for trace messages
	LogLevelTrace = 16
)

const (
	// LogDateFormat is the format for the timestamp of log entries
	LogDateFormat string = "2006-01-02 15:04:05.000 MST"
)

var logFatal = Logger.Fatal

// Logger is a custom log writer that adds a timestamp to each log entry
type Logger struct {
	Level  LogLevel
	Output *os.File
	m      *sync.Mutex
}

// Options are options for the Logger
type Options struct {
	Level  LogLevel
	Output *os.File
}

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "Debug"
	case LogLevelError:
		return "Error"
	case LogLevelFatal:
		return "Fatal"
	case LogLevelInfo:
		return "Info"
	case LogLevelTrace:
		return "Trace"
	case LogLevelVerbose:
		return "Verbose"
	case LogLevelWarn:
		return "Warn"
	default:
		return "Unknown"
	}
}

// buildMessage builds a log message with a prefix and args passed to it
// The arguments are separated by a space
func buildMessage(l LogLevel, a ...any) string {

	var message string

	prefix := fmt.Sprintf(time.Now().Format(LogDateFormat)) + " " + l.String()[0:1] + " "
	prefixLength := len(prefix)

	for _, v := range a {
		lines := strings.Split(v.(string), "\n")
		for ix, line := range lines {
			if line != "" {
				if ix == 0 {
					message += prefix + line + "\n"
				} else {
					message += strings.Repeat(" ", prefixLength) + line + "\n"
				}
			}
		}
	}

	message = strings.TrimRight(message, " ")
	// if message[len(message)-1] != ' ' {
	// 	message += " "
	// }

	return message
}

// New creates a new Logger
func New() *Logger {
	return &Logger{Level: LogLevelInfo, Output: os.Stderr, m: &sync.Mutex{}}
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(opts Options) *Logger {

	if opts.Level == 0 {
		opts.Level = LogLevelInfo
	}

	if opts.Output == nil {
		opts.Output = os.Stderr
	}

	return &Logger{Level: opts.Level, Output: opts.Output, m: &sync.Mutex{}}
}

// SetLevel sets the log level
func (writer *Logger) SetLevel(level LogLevel) {

	writer.m.Lock()
	writer.Level = level
	writer.m.Unlock()
	writer.Infof("log level set to %s\n", level.String())

	return
}

// SetOutput sets the output file for the logger
func (writer *Logger) SetOutput(file *os.File) {

	writer.m.Lock()

	if file == nil {
		writer.Output = os.Stderr
	} else {
		writer.Output = file
	}

	writer.m.Unlock()
	writer.Debugf("output set to %s\n", file.Name())

	return
}

// Write writes a log entry to an output file (default: os.Stderr)
func (writer Logger) Write(bytes []byte) (int, error) {

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
func (writer Logger) Debug(a ...any) {

	if writer.Level >= LogLevelDebug {
		message := buildMessage(LogLevelDebug, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Debugf logs a debug message with a format string
func (writer Logger) Debugf(format string, a ...any) {

	if writer.Level >= LogLevelDebug {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelDebug, msg)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Error logs an error message
func (writer Logger) Error(a ...any) {

	if writer.Level >= LogLevelError {
		message := buildMessage(LogLevelError, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Errorf logs an error message with a format string
func (writer Logger) Errorf(format string, a ...any) {

	if writer.Level >= LogLevelError {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelError, msg)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Fatal logs a fatal message
func (writer Logger) Fatal(a ...any) {

	message := buildMessage(LogLevelFatal, a...)
	_, err := writer.Write([]byte(message))

	if err != nil {
		panic(err)
	}

	os.Exit(1)
}

// Fatalf logs a fatal message with a format string
func (writer Logger) Fatalf(format string, a ...any) {

	msg := fmt.Sprintf(format, a...)
	message := buildMessage(LogLevelFatal, msg)
	_, err := writer.Write([]byte(message))

	if err != nil {
		panic(err)
	}

	os.Exit(1)
}

// Info logs an info message
func (writer Logger) Info(a ...any) {

	if writer.Level >= LogLevelInfo {
		message := buildMessage(LogLevelInfo, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Infof logs an info message with a format string
func (writer Logger) Infof(format string, a ...any) {

	if writer.Level >= LogLevelInfo {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelInfo, msg)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Trace logs a trace message
func (writer Logger) Trace(a ...any) {

	if writer.Level >= LogLevelTrace {
		message := buildMessage(LogLevelTrace, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Tracef logs a warning message with a format string
func (writer Logger) Tracef(format string, a ...any) {

	if writer.Level >= LogLevelTrace {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelTrace, msg)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Verbose logs a verbose message
func (writer Logger) Verbose(a ...any) {

	if writer.Level >= LogLevelVerbose {
		message := buildMessage(LogLevelVerbose, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Verbosef logs a verbose message with a format string
func (writer Logger) Verbosef(format string, a ...any) {

	if writer.Level >= LogLevelVerbose {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelVerbose, msg)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Warn logs a warning message
func (writer Logger) Warn(a ...any) {

	if writer.Level >= LogLevelWarn {
		message := buildMessage(LogLevelWarn, a...)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}

// Warnf logs a warning message with a format string
func (writer Logger) Warnf(format string, a ...any) {

	if writer.Level >= LogLevelWarn {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelWarn, msg)
		_, err := writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}

	return
}
