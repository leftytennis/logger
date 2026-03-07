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
	LogLevelFatal LogLevel = -16
	// LogLevelError is the log level for error messages
	LogLevelError LogLevel = -8
	// LogLevelWarn is the log level for warning messages
	LogLevelWarn LogLevel = -4
	// LogLevelInfo is the log level for info messages
	LogLevelInfo LogLevel = 0
	// LogLevelVerbose is the log level for verbose messages
	LogLevelVerbose LogLevel = 4
	// LogLevelDebug is the log level for debug messages
	LogLevelDebug LogLevel = 8
	// LogLevelTrace is the log level for trace messages
	LogLevelTrace LogLevel = 16
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
	exitFunc   func(int)
	m          sync.RWMutex
}

// Options are options for the Logger
type Options struct {
	Level      LogLevel
	LevelCount int
	Output     io.Writer
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
// Multiple args are space-separated (like fmt.Println). Newlines in the
// resulting string get continuation-line padding.
func buildMessage(l LogLevel, a ...any) string {

	prefix := time.Now().Format(LogDateFormat) + " " + l.String()[0:1] + " "
	padding := strings.Repeat(" ", len(prefix))

	combined := fmt.Sprint(a...)
	if combined == "" {
		return prefix + "\n"
	}

	lines := strings.Split(combined, "\n")

	var b strings.Builder
	for i, line := range lines {
		if line != "" {
			if i == 0 {
				b.WriteString(prefix)
			} else {
				b.WriteString(padding)
			}
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	result := b.String()
	if result == "" {
		return prefix + "\n"
	}
	return result
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
	case "trace", "trace1":
		logLevel = LogLevelTrace
		levelCount = 1
	case "trace2":
		logLevel = LogLevelTrace
		levelCount = 2
	case "trace3":
		logLevel = LogLevelTrace
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
		LevelCount: 1,
		Output:     os.Stderr,
	})
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(ctxParent context.Context, opts Options) *Logger {

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
		exitFunc:   os.Exit,
	}

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

// levelAllowed checks if the given level should be logged
func (l *Logger) levelAllowed(level LogLevel) bool {
	l.m.RLock()
	allowed := l.Level >= level
	l.m.RUnlock()
	return allowed
}

// subLevelAllowed checks if a sub-level (LevelCount) log should be logged
func (l *Logger) subLevelAllowed(level LogLevel, requiredCount int) bool {
	l.m.RLock()
	allowed := l.Level > level || (l.Level == level && l.LevelCount >= requiredCount)
	l.m.RUnlock()
	return allowed
}

// log is the internal method that handles level checking, message building, and writing
func (l *Logger) log(level LogLevel, a ...any) {
	if l.levelAllowed(level) {
		message := buildMessage(level, a...)
		l.Write([]byte(message))
	}
}

// logf is the internal method for formatted log messages
func (l *Logger) logf(level LogLevel, format string, a ...any) {
	if l.levelAllowed(level) {
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
	if l.subLevelAllowed(LogLevelDebug, 2) {
		message := buildMessage(LogLevelDebug, a...)
		l.Write([]byte(message))
	}
}

// Debug3 logs a debug message when LevelCount is >= 3 (-ddd)
func (l *Logger) Debug3(a ...any) {
	if l.subLevelAllowed(LogLevelDebug, 3) {
		message := buildMessage(LogLevelDebug, a...)
		l.Write([]byte(message))
	}
}

// Debugf logs a debug message with a format string
func (l *Logger) Debugf(format string, a ...any) {
	l.logf(LogLevelDebug, format, a...)
}

// Debugf2 logs a debug message with a format string when LevelCount >= 2 (-dd)
func (l *Logger) Debugf2(format string, a ...any) {
	if l.subLevelAllowed(LogLevelDebug, 2) {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelDebug, msg)
		l.Write([]byte(message))
	}
}

// Debugf3 logs a debug message with a format string when LevelCount >= 3 (-ddd)
func (l *Logger) Debugf3(format string, a ...any) {
	if l.subLevelAllowed(LogLevelDebug, 3) {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelDebug, msg)
		l.Write([]byte(message))
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
	l.exitFunc(1)
}

// Fatalf logs a fatal message with a format string
func (l *Logger) Fatalf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	message := buildMessage(LogLevelFatal, msg)
	l.Write([]byte(message))
	l.exitFunc(1)
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
	if l.subLevelAllowed(LogLevelVerbose, 2) {
		message := buildMessage(LogLevelVerbose, a...)
		l.Write([]byte(message))
	}
}

// Verbose3 logs a verbose message when LevelCount >= 3 (i.e., -vvv)
func (l *Logger) Verbose3(a ...any) {
	if l.subLevelAllowed(LogLevelVerbose, 3) {
		message := buildMessage(LogLevelVerbose, a...)
		l.Write([]byte(message))
	}
}

// Verbosef logs a verbose message with a format string
func (l *Logger) Verbosef(format string, a ...any) {
	l.logf(LogLevelVerbose, format, a...)
}

// Verbosef2 logs a verbose message with a format string when LevelCount >= 2 (i.e., -vv)
func (l *Logger) Verbosef2(format string, a ...any) {
	if l.subLevelAllowed(LogLevelVerbose, 2) {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelVerbose, msg)
		l.Write([]byte(message))
	}
}

// Verbosef3 logs a verbose message with a format string when LevelCount >= 3 (i.e., -vvv)
func (l *Logger) Verbosef3(format string, a ...any) {
	if l.subLevelAllowed(LogLevelVerbose, 3) {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(LogLevelVerbose, msg)
		l.Write([]byte(message))
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
