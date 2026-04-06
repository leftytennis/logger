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

// Level is a type for log levels
type Level int

// loggerKey is a type for the logger key stored in a context
type loggerKeyType struct{}

const (
	// LevelFatal is the highest log level and will log fatal messages
	LevelFatal Level = -16
	// LevelError is the log level for error messages
	LevelError Level = -8
	// LevelWarn is the log level for warning messages
	LevelWarn Level = -4
	// LevelInfo is the log level for info messages
	LevelInfo Level = 0
	// LevelVerbose is the log level for verbose messages
	LevelVerbose Level = 4
	// LevelDebug is the log level for debug messages
	LevelDebug Level = 8
	// LevelTrace is the log level for trace messages
	LevelTrace Level = 16
)

const (
	// DateFormat is the format for the timestamp of log entries
	DateFormat string = "2006-01-02 15:04:05.000 MST"
)

var (
	loggerKey = loggerKeyType{}
)

// Logger is a custom log writer that adds a timestamp to each log entry
type Logger struct {
	level      Level
	levelCount int
	output     io.Writer
	exitFunc   func(int)
	m          sync.RWMutex
}

// Options are options for the Logger
type Options struct {
	Level      Level
	LevelCount int
	Output     io.Writer
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "Debug"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "Fatal"
	case LevelInfo:
		return "Info"
	case LevelTrace:
		return "Trace"
	case LevelVerbose:
		return "Verbose"
	case LevelWarn:
		return "Warn"
	default:
		return "Unknown"
	}
}

// buildMessage builds a log message with a prefix and args passed to it.
// Multiple args are space-separated (like fmt.Println). Newlines in the
// resulting string get continuation-line padding.
func buildMessage(l Level, a ...any) string {

	prefix := time.Now().Format(DateFormat) + " " + l.String()[0:1] + " "
	padding := strings.Repeat(" ", len(prefix))

	combined := strings.TrimRight(fmt.Sprintln(a...), "\n")
	if combined == "" {
		return prefix + "\n"
	}

	lines := strings.Split(combined, "\n")

	var b strings.Builder
	b.Grow(len(prefix) + len(combined) + len(lines)*len(padding))
	for i, line := range lines {
		// Skip only a trailing empty element produced by a terminal newline
		if i == len(lines)-1 && line == "" {
			continue
		}
		if i == 0 {
			b.WriteString(prefix)
		} else {
			b.WriteString(padding)
		}
		b.WriteString(line)
		b.WriteByte('\n')
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

// ParseLevel returns the Level and level count from a string (i.e. "debug", "debug1"...).
// Returns an error for unknown inputs.
func ParseLevel(s string) (Level, int, error) {

	var logLevel Level
	var levelCount int = 1

	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info":
		logLevel = LevelInfo
	case "verbose", "verbose1":
		logLevel = LevelVerbose
		levelCount = 1
	case "verbose2":
		logLevel = LevelVerbose
		levelCount = 2
	case "verbose3":
		logLevel = LevelVerbose
		levelCount = 3
	case "debug", "debug1":
		logLevel = LevelDebug
		levelCount = 1
	case "debug2":
		logLevel = LevelDebug
		levelCount = 2
	case "debug3":
		logLevel = LevelDebug
		levelCount = 3
	case "trace", "trace1":
		logLevel = LevelTrace
		levelCount = 1
	case "trace2":
		logLevel = LevelTrace
		levelCount = 2
	case "trace3":
		logLevel = LevelTrace
		levelCount = 3
	case "warn", "warning":
		logLevel = LevelWarn
	case "error":
		logLevel = LevelError
	case "fatal":
		logLevel = LevelFatal
	default:
		return LevelInfo, 1, fmt.Errorf("unknown log level: %q", s)
	}

	return logLevel, levelCount, nil
}

// New creates a new Logger with default settings (Info level, stderr output)
func New() *Logger {
	return NewWithOptions(Options{
		Level:      LevelInfo,
		LevelCount: 1,
		Output:     os.Stderr,
	})
}

// NewWithOptions creates a new Logger with options
func NewWithOptions(opts Options) *Logger {

	if opts.LevelCount == 0 {
		opts.LevelCount = 1
	}

	if opts.Output == nil {
		opts.Output = os.Stderr
	}

	// Create new logger
	log := &Logger{
		level:      opts.Level,
		levelCount: opts.LevelCount,
		output:     opts.Output,
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

// GetLevel returns the current log level
func (l *Logger) GetLevel() Level {
	l.m.RLock()
	level := l.level
	l.m.RUnlock()
	return level
}

// GetLevelCount returns the current level count
func (l *Logger) GetLevelCount() int {
	l.m.RLock()
	count := l.levelCount
	l.m.RUnlock()
	return count
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level Level, levelCount int) {

	l.m.Lock()
	l.level = level
	l.levelCount = levelCount
	l.m.Unlock()
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(w io.Writer) {

	l.m.Lock()

	if w == nil {
		l.output = os.Stderr
	} else {
		l.output = w
	}

	l.m.Unlock()

}

// Write writes a log entry to an output writer (default: os.Stderr)
func (l *Logger) Write(bytes []byte) (int, error) {

	l.m.Lock()
	defer l.m.Unlock()

	if l.output == nil {
		l.output = os.Stderr
	}

	if len(bytes) > 0 && bytes[len(bytes)-1] != '\n' {
		// Force a new allocation to avoid mutating the caller's backing array
		bytes = append(bytes[:len(bytes):len(bytes)], '\n')
	}

	return l.output.Write(bytes)
}

// Enabled reports whether the logger would emit a log entry at the given
// level and count. Useful for guarding expensive log-arg construction:
//
//	if log.Enabled(logger.LevelDebug, 2) {
//	    log.Debug2(expensiveString())
//	}
func (l *Logger) Enabled(level Level, count int) bool {
	if count <= 1 {
		return l.levelAllowed(level)
	}
	return l.subLevelAllowed(level, count)
}

// levelAllowed checks if the given level should be logged
func (l *Logger) levelAllowed(level Level) bool {
	l.m.RLock()
	allowed := l.level >= level
	l.m.RUnlock()
	return allowed
}

// subLevelAllowed checks if a sub-level (levelCount) log should be logged.
// When the logger is set to a higher level (e.g., Trace), all sub-levels
// of lower levels (Debug2, Verbose3, etc.) are automatically included,
// since l.level > level is true for those cases.
func (l *Logger) subLevelAllowed(level Level, requiredCount int) bool {
	l.m.RLock()
	allowed := l.level > level || (l.level == level && l.levelCount >= requiredCount)
	l.m.RUnlock()
	return allowed
}

// log is the internal method that handles level checking, message building, and writing
func (l *Logger) log(level Level, a ...any) {
	if l.levelAllowed(level) {
		message := buildMessage(level, a...)
		l.Write([]byte(message))
	}
}

// logf is the internal method for formatted log messages
func (l *Logger) logf(level Level, format string, a ...any) {
	if l.levelAllowed(level) {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(level, msg)
		l.Write([]byte(message))
	}
}

// logSub is the internal method for sub-level log messages (e.g., Debug2, Verbose3)
func (l *Logger) logSub(level Level, requiredCount int, a ...any) {
	if l.subLevelAllowed(level, requiredCount) {
		message := buildMessage(level, a...)
		l.Write([]byte(message))
	}
}

// logSubf is the internal method for formatted sub-level log messages
func (l *Logger) logSubf(level Level, requiredCount int, format string, a ...any) {
	if l.subLevelAllowed(level, requiredCount) {
		msg := fmt.Sprintf(format, a...)
		message := buildMessage(level, msg)
		l.Write([]byte(message))
	}
}

// Debug logs a debug message
func (l *Logger) Debug(a ...any) {
	l.log(LevelDebug, a...)
}

// Debug2 logs a debug message when LevelCount is >= 2 (-dd)
func (l *Logger) Debug2(a ...any) {
	l.logSub(LevelDebug, 2, a...)
}

// Debug3 logs a debug message when LevelCount is >= 3 (-ddd)
func (l *Logger) Debug3(a ...any) {
	l.logSub(LevelDebug, 3, a...)
}

// Debugf logs a debug message with a format string
func (l *Logger) Debugf(format string, a ...any) {
	l.logf(LevelDebug, format, a...)
}

// Debugf2 logs a debug message with a format string when LevelCount >= 2 (-dd)
func (l *Logger) Debugf2(format string, a ...any) {
	l.logSubf(LevelDebug, 2, format, a...)
}

// Debugf3 logs a debug message with a format string when LevelCount >= 3 (-ddd)
func (l *Logger) Debugf3(format string, a ...any) {
	l.logSubf(LevelDebug, 3, format, a...)
}

// Error logs an error message
func (l *Logger) Error(a ...any) {
	l.log(LevelError, a...)
}

// Errorf logs an error message with a format string
func (l *Logger) Errorf(format string, a ...any) {
	l.logf(LevelError, format, a...)
}

// Fatal logs a fatal message
func (l *Logger) Fatal(a ...any) {
	message := buildMessage(LevelFatal, a...)
	l.Write([]byte(message))
	l.exitFunc(1)
}

// Fatalf logs a fatal message with a format string
func (l *Logger) Fatalf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	message := buildMessage(LevelFatal, msg)
	l.Write([]byte(message))
	l.exitFunc(1)
}

// Info logs an info message
func (l *Logger) Info(a ...any) {
	l.log(LevelInfo, a...)
}

// Infof logs an info message with a format string
func (l *Logger) Infof(format string, a ...any) {
	l.logf(LevelInfo, format, a...)
}

// Trace logs a trace message
func (l *Logger) Trace(a ...any) {
	l.log(LevelTrace, a...)
}

// Trace2 logs a trace message when LevelCount is >= 2
func (l *Logger) Trace2(a ...any) {
	l.logSub(LevelTrace, 2, a...)
}

// Trace3 logs a trace message when LevelCount is >= 3
func (l *Logger) Trace3(a ...any) {
	l.logSub(LevelTrace, 3, a...)
}

// Tracef logs a trace message with a format string
func (l *Logger) Tracef(format string, a ...any) {
	l.logf(LevelTrace, format, a...)
}

// Tracef2 logs a trace message with a format string when LevelCount >= 2
func (l *Logger) Tracef2(format string, a ...any) {
	l.logSubf(LevelTrace, 2, format, a...)
}

// Tracef3 logs a trace message with a format string when LevelCount >= 3
func (l *Logger) Tracef3(format string, a ...any) {
	l.logSubf(LevelTrace, 3, format, a...)
}

// Verbose logs a verbose message
func (l *Logger) Verbose(a ...any) {
	l.log(LevelVerbose, a...)
}

// Verbose2 logs a verbose message when LevelCount >= 2 (i.e., -vv)
func (l *Logger) Verbose2(a ...any) {
	l.logSub(LevelVerbose, 2, a...)
}

// Verbose3 logs a verbose message when LevelCount >= 3 (i.e., -vvv)
func (l *Logger) Verbose3(a ...any) {
	l.logSub(LevelVerbose, 3, a...)
}

// Verbosef logs a verbose message with a format string
func (l *Logger) Verbosef(format string, a ...any) {
	l.logf(LevelVerbose, format, a...)
}

// Verbosef2 logs a verbose message with a format string when LevelCount >= 2 (i.e., -vv)
func (l *Logger) Verbosef2(format string, a ...any) {
	l.logSubf(LevelVerbose, 2, format, a...)
}

// Verbosef3 logs a verbose message with a format string when LevelCount >= 3 (i.e., -vvv)
func (l *Logger) Verbosef3(format string, a ...any) {
	l.logSubf(LevelVerbose, 3, format, a...)
}

// Warn logs a warning message
func (l *Logger) Warn(a ...any) {
	l.log(LevelWarn, a...)
}

// Warnf logs a warning message with a format string
func (l *Logger) Warnf(format string, a ...any) {
	l.logf(LevelWarn, format, a...)
}
