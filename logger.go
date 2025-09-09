// Package logger provides a custom log writer that adds a timestamp to each log entry.
package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	// "github.com/leftytennis/logger"
)

// LogLevel is a type for log levels
type LogLevel int

// loggerKey is a type for the logger key stored in a context
type loggerKeyType struct{}

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

var (
	ctx                  context.Context
	loggerKey            = loggerKeyType{}
	logLevelCount        int
)

// var logFatal = Logger.Fatal

// Logger is a custom log writer that adds a timestamp to each log entry
type Logger struct {
	Context    context.Context
	Level      LogLevel
	LevelCount int
	Output     *os.File
	m          *sync.Mutex
}

// Options are options for the Logger
type Options struct {
	Level      LogLevel
	LevelCount int
	Output     *os.File
}

func init() {
	ctx = context.TODO()
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
		return "Info" // LogLevelVerbose is treated as LogLevelInfo
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

	if len(message) == 0 {
		message += prefix + "\n"
	}

	return message
}

// FromContext returns a pointer to the logger from a context
func FromContext(ctxParent context.Context) *Logger {

	// return logger pointer if it's in the context
	if value, exists := ctxParent.Value(loggerKey).(*Logger); exists {
		return value
	}

	// return nil if logger is not found
	return nil
}

// New creates a new Logger
func New(ctx context.Context) *Logger {
	return NewWithOptions(ctx, Options{
		Level:  LogLevelInfo,
		LevelCount: 1,
		Output: os.Stderr,
	})
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(ctxParent context.Context, opts Options) *Logger {

	if opts.Level == 0 {
		opts.Level = LogLevelInfo
	}

	if opts.LevelCount == 0 {
		opts.LevelCount = 1
	}

	if opts.Output == nil {
		opts.Output = os.Stderr
	}

	logLevelCount = opts.LevelCount
	
	// Create new logger
	log := &Logger{
		Level:  opts.Level,
		LevelCount: opts.LevelCount,
		Output: opts.Output,
		m:      &sync.Mutex{},
	}

	// create new context and store logger value
	ctx = WithContext(ctxParent, log)

	return log
}

// WithContext returns a copy of the logger with the context set to ctx
func WithContext(ctxParent context.Context, log *Logger) context.Context {

	var value *Logger
	var exists bool

	// see if context already has the logger
	if value, exists = ctxParent.Value(loggerKey).(*Logger); exists {
		if value == log {
			return ctxParent
		}
	}

	// create a new context that has the logger value
	ctx := context.WithValue(ctxParent, loggerKey, log)

	return ctx
}

// SetLevel sets the log level
func (writer *Logger) SetLevel(level LogLevel, levelCount int) {

	writer.m.Lock()
	writer.Level = level
	writer.LevelCount = levelCount
	writer.m.Unlock()
	writer.Infof("log level set to %s\n", level.String())

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

}

// Write writes a log entry to an output file (default: os.Stdout)
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

}

// Debug2 logs a debug message when logLevelCount is >= 2 (-dd)
func (writer Logger) Debug2(a ...any) {

	if writer.Level > LogLevelDebug || ((writer.Level == LogLevelDebug) && (logLevelCount >= 2)) {
		writer.Debug(a)
	}

}

// Debug3 logs a debug message when logLevelCount is >= 3 (-ddd)
func (writer Logger) Debug3(a ...any) {

	if writer.Level > LogLevelDebug || ((writer.Level == LogLevelDebug) && (logLevelCount >= 3)) {
		writer.Debug(a...)
	}

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

}

// Debugf2 logs a debug message with a format string when logLevelCount >= 2 (-dd)
func (writer Logger) Debugf2(format string, a ...any) {

	if writer.Level > LogLevelDebug || ((writer.Level == LogLevelDebug) && (logLevelCount >= 2)) {
		writer.Debugf(format, a...)
	}

}

// Debugf3 logs a debug message with a format string when logLevelCount >= 3 (-ddd)
func (writer Logger) Debugf3(format string, a ...any) {

	if writer.Level > LogLevelDebug || ((writer.Level == LogLevelDebug) && (logLevelCount >= 3)) {
		writer.Debugf(format, a...)
	}

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

}

// Verbose2 logs a verbose message when logLevelCount >= 2 (i.e., -vv)
func (writer Logger) Verbose2(a ...any) {

	if writer.Level > LogLevelVerbose || ((writer.Level == LogLevelVerbose) && (logLevelCount >= 2)) {
		writer.Verbose(a...)
	}

}

// Verbose3 logs a verbose message when logLevelCount >= 3 (i.e., -vvv)
func (writer Logger) Verbose3(a ...any) {

	if writer.Level > LogLevelVerbose || ((writer.Level == LogLevelVerbose) && (logLevelCount >= 3)) {
		writer.Verbose(a...)
	}

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

}

// Verbosef2 logs a verbose message with a format string when logLevelCount >= 2 (i.e., -vv)
func (writer Logger) Verbosef2(format string, a ...any) {

	if writer.Level > LogLevelVerbose || ((writer.Level == LogLevelVerbose) && (logLevelCount >= 2)) {
		writer.Verbosef(format, a...)
	}

}

// Verbosef3 logs a verbose message with a format string when logLevelCount >= 3 (i.e., -vvv)
func (writer Logger) Verbosef3(format string, a ...any) {

	if writer.Level >= LogLevelVerbose || ((writer.Level == LogLevelVerbose) && (logLevelCount >= 3)) {
		writer.Verbosef(format, a...)
	}

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

}
