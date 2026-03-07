// Package logger provides a custom log writer that adds a timestamp to each log entry.
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
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
	loggerKey = loggerKeyType{}
)

// Logger is a custom log writer that adds a timestamp to each log entry
type Logger struct {
	Level      LogLevel
	LevelCount int
	Output     io.Writer
	m          sync.Mutex
}

// Options are options for the Logger
type Options struct {
	Level      LogLevel
	LevelCount int
	Output     io.Writer
	LevelSet   bool // explicitly indicate that Level was set
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

	prefix := time.Now().Format(LogDateFormat) + " " + l.String()[0:1] + " "
	prefixLength := len(prefix)

	for _, v := range a {
		lines := strings.Split(fmt.Sprint(v), "\n")
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

// GetLoggerValuesFromString returns the LogLevel and level count from a string (i.e. "debug", "debug1"...)
func GetLoggerValuesFromString(levelStr string) (LogLevel, int) {

	var logLevel LogLevel
	var levelCount int = 1

	switch strings.ToLower(levelStr) {
	case "info":
		logLevel = LogLevelInfo
	case "verbose", "verbose1":
		logLevel = LogLevelVerbose
		levelCount = 1
	case "verbose2":
		logLevel = LogLevelVerbose
		levelCount = 2
	case "verbose3":
		logLevel = LogLevelVerbose
		levelCount = 3
	case "debug", "debug1":
		logLevel = LogLevelDebug
		levelCount = 1
	case "debug2":
		logLevel = LogLevelDebug
		levelCount = 2
	case "debug3":
		logLevel = LogLevelDebug
		levelCount = 3
	case "warn", "warning":
		logLevel = LogLevelWarn
	case "error":
		logLevel = LogLevelError
	case "fatal":
		logLevel = LogLevelFatal
	default:
		logLevel = LogLevelInfo
	}

	return logLevel, levelCount
}

// New creates a new Logger
func New(ctx context.Context) *Logger {
	return NewWithOptions(ctx, Options{
		Level:      LogLevelInfo,
		LevelSet:   true,
		LevelCount: 1,
		Output:     os.Stderr,
	})
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(ctxParent context.Context, opts Options) *Logger {

	if !opts.LevelSet && opts.Level == 0 {
		opts.Level = LogLevelInfo
	}

	if opts.LevelCount == 0 {
		opts.LevelCount = 1
	}

	if opts.Output == nil {
		opts.Output = os.Stderr
	}

	// Create new logger
	log := &Logger{
		Level:      opts.Level,
		LevelCount: opts.LevelCount,
		Output:     opts.Output,
	}

	// create new context and store logger value
	WithContext(ctxParent, log)

	return log
}

// WithContext returns a copy of the logger with the context set to ctx
func WithContext(ctxParent context.Context, log *Logger) context.Context {

	// see if context already has the logger
	if value, exists := ctxParent.Value(loggerKey).(*Logger); exists {
		if value == log {
			return ctxParent
		}
	}

	// create a new context that has the logger value
	ctx := context.WithValue(ctxParent, loggerKey, log)

	return ctx
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level LogLevel, levelCount int) {

	l.m.Lock()
	l.Level = level
	l.LevelCount = levelCount
	l.m.Unlock()
	l.Infof("log level set to %s\n", level.String())

}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(w io.Writer) {

	l.m.Lock()

	if w == nil {
		l.Output = os.Stderr
	} else {
		l.Output = w
	}

	l.m.Unlock()

}

// Write writes a log entry to an output writer (default: os.Stderr)
func (l *Logger) Write(bytes []byte) (int, error) {

	l.m.Lock()
	defer l.m.Unlock()

	if l.Output == nil {
		l.Output = os.Stderr
	}

	if len(bytes) > 0 && bytes[len(bytes)-1] != '\n' {
		bytes = append(bytes, '\n')
	}

	return l.Output.Write(bytes)
}

// log is the internal method that handles level checking, message building, and writing
func (l *Logger) log(level LogLevel, a ...any) {
	if l.Level >= level {
		message := buildMessage(level, a...)
		l.Write([]byte(message))
	}
}

// logf is the internal method for formatted log messages
func (l *Logger) logf(level LogLevel, format string, a ...any) {
	if l.Level >= level {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(level, msg)
		l.Write([]byte(message))
	}
}

// Debug logs a debug message
func (l *Logger) Debug(a ...any) {
	l.log(LogLevelDebug, a...)
}

// Debug2 logs a debug message when LevelCount is >= 2 (-dd)
func (l *Logger) Debug2(a ...any) {
	if l.Level > LogLevelDebug || (l.Level == LogLevelDebug && l.LevelCount >= 2) {
		l.Debug(a...)
	}
}

// Debug3 logs a debug message when LevelCount is >= 3 (-ddd)
func (l *Logger) Debug3(a ...any) {
	if l.Level > LogLevelDebug || (l.Level == LogLevelDebug && l.LevelCount >= 3) {
		l.Debug(a...)
	}
}

// Debugf logs a debug message with a format string
func (l *Logger) Debugf(format string, a ...any) {
	l.logf(LogLevelDebug, format, a...)
}

// Debugf2 logs a debug message with a format string when LevelCount >= 2 (-dd)
func (l *Logger) Debugf2(format string, a ...any) {
	if l.Level > LogLevelDebug || (l.Level == LogLevelDebug && l.LevelCount >= 2) {
		l.Debugf(format, a...)
	}
}

// Debugf3 logs a debug message with a format string when LevelCount >= 3 (-ddd)
func (l *Logger) Debugf3(format string, a ...any) {
	if l.Level > LogLevelDebug || (l.Level == LogLevelDebug && l.LevelCount >= 3) {
		l.Debugf(format, a...)
	}
}

// Error logs an error message
func (l *Logger) Error(a ...any) {
	l.log(LogLevelError, a...)
}

// Errorf logs an error message with a format string
func (l *Logger) Errorf(format string, a ...any) {
	l.logf(LogLevelError, format, a...)
}

// Fatal logs a fatal message
func (l *Logger) Fatal(a ...any) {
	message := buildMessage(LogLevelFatal, a...)
	l.Write([]byte(message))
	os.Exit(1)
}

// Fatalf logs a fatal message with a format string
func (l *Logger) Fatalf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	message := buildMessage(LogLevelFatal, msg)
	l.Write([]byte(message))
	os.Exit(1)
}

// Info logs an info message
func (l *Logger) Info(a ...any) {
	l.log(LogLevelInfo, a...)
}

// Infof logs an info message with a format string
func (l *Logger) Infof(format string, a ...any) {
	l.logf(LogLevelInfo, format, a...)
}

// Trace logs a trace message
func (l *Logger) Trace(a ...any) {
	l.log(LogLevelTrace, a...)
}

// Tracef logs a trace message with a format string
func (l *Logger) Tracef(format string, a ...any) {
	l.logf(LogLevelTrace, format, a...)
}

// Verbose logs a verbose message
func (l *Logger) Verbose(a ...any) {
	l.log(LogLevelVerbose, a...)
}

// Verbose2 logs a verbose message when LevelCount >= 2 (i.e., -vv)
func (l *Logger) Verbose2(a ...any) {
	if l.Level > LogLevelVerbose || (l.Level == LogLevelVerbose && l.LevelCount >= 2) {
		l.Verbose(a...)
	}
}

// Verbose3 logs a verbose message when LevelCount >= 3 (i.e., -vvv)
func (l *Logger) Verbose3(a ...any) {
	if l.Level > LogLevelVerbose || (l.Level == LogLevelVerbose && l.LevelCount >= 3) {
		l.Verbose(a...)
	}
}

// Verbosef logs a verbose message with a format string
func (l *Logger) Verbosef(format string, a ...any) {
	l.logf(LogLevelVerbose, format, a...)
}

// Verbosef2 logs a verbose message with a format string when LevelCount >= 2 (i.e., -vv)
func (l *Logger) Verbosef2(format string, a ...any) {
	if l.Level > LogLevelVerbose || (l.Level == LogLevelVerbose && l.LevelCount >= 2) {
		l.Verbosef(format, a...)
	}
}

// Verbosef3 logs a verbose message with a format string when LevelCount >= 3 (i.e., -vvv)
func (l *Logger) Verbosef3(format string, a ...any) {
	if l.Level > LogLevelVerbose || (l.Level == LogLevelVerbose && l.LevelCount >= 3) {
		l.Verbosef(format, a...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(a ...any) {
	l.log(LogLevelWarn, a...)
}

// Warnf logs a warning message with a format string
func (l *Logger) Warnf(format string, a ...any) {
	l.logf(LogLevelWarn, format, a...)
}
