package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {

	ctx := context.TODO()
	opts := Options{}

	logr := NewWithOptions(ctx, opts)

	if logr == nil {
		t.Fatal("Expected a valid Logger instance, got nil")
	}

	ctx = WithContext(ctx, logr)

	logrValue := FromContext(ctx)
	if logrValue == nil {
		t.Fatal("Expected a valid Logger instance from context, got nil")
	}

	if logr.Level != LogLevelInfo {
		t.Errorf("Expected default log level to be LogLevelInfo, got %v", logr.Level)
	}

	if logr.Output != os.Stderr {
		t.Errorf("Expected default output to be os.Stderr, got %v", logr.Output)
	}

}

func TestNewWithOptions(t *testing.T) {

	ctx := context.TODO()
	opts := Options{
		Level:    LogLevelDebug,
		LevelSet: true,
		Output:   os.Stdout,
	}

	logr := NewWithOptions(ctx, opts)

	if logr.Level != LogLevelDebug {
		t.Errorf("Expected log level to be LogLevelDebug, got %v", logr.Level)
	}

	if logr.Output != os.Stdout {
		t.Errorf("Expected output to be os.Stdout, got %v", logr.Output)
	}

}

func TestSetLevel(t *testing.T) {

	logr := New(context.Background())
	logr.SetLevel(LogLevelError, 1)

	if logr.Level != LogLevelError {
		t.Errorf("Expected log level to be LogLevelError, got %v", logr.Level)
	}

}

func TestSetOutput(t *testing.T) {

	logr := New(context.Background())
	var buf bytes.Buffer

	logr.SetOutput(&buf)

	logr.Info("test output")

	if buf.Len() == 0 {
		t.Error("Expected output to be written to buffer")
	}

}

func TestDebug(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, Output: &buf})

	logr.Debug("Debug message")
	output := buf.String()
	if !strings.Contains(output, "D Debug message") {
		t.Errorf("Expected debug message in output, got %q", output)
	}

	buf.Reset()
	logr.Debug("Debug message\nwith newline")
	output = buf.String()
	if !strings.Contains(output, "D Debug message") || !strings.Contains(output, "with newline") {
		t.Errorf("Expected multiline debug message in output, got %q", output)
	}

}

func TestDebugf(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, Output: &buf})

	logr.Debugf("Debug message %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "D Debug message formatted") {
		t.Errorf("Expected formatted debug message in output, got %q", output)
	}

}

func TestDebugLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, LevelCount: 1, Output: &buf})

	logr.Trace("trace message")
	if buf.Len() != 0 {
		t.Error("Trace should not be logged at Debug level")
	}

	logr.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("Debug should be logged at Debug level")
	}

	buf.Reset()
	logr.Info("info message")
	if buf.Len() == 0 {
		t.Error("Info should be logged at Debug level")
	}

	buf.Reset()
	logr.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("Warn should be logged at Debug level")
	}

	buf.Reset()
	logr.Error("error message")
	if buf.Len() == 0 {
		t.Error("Error should be logged at Debug level")
	}

}

func TestError(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelError, LevelSet: true, Output: &buf})

	logr.Error("Error message")
	output := buf.String()
	if !strings.Contains(output, "E Error message") {
		t.Errorf("Expected error message in output, got %q", output)
	}

}

func TestErrorf(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelError, LevelSet: true, Output: &buf})

	logr.Errorf("Error message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "E Error message: formatted") {
		t.Errorf("Expected formatted error message in output, got %q", output)
	}

}

func TestErrorLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelError, LevelSet: true, Output: &buf})

	logr.Trace("trace message")
	logr.Debug("debug message")
	logr.Verbose("verbose message")
	logr.Info("info message")
	logr.Warn("warn message")
	if buf.Len() != 0 {
		t.Error("Nothing below Error should be logged at Error level")
	}

	logr.Error("error message")
	if buf.Len() == 0 {
		t.Error("Error should be logged at Error level")
	}

}

func TestInfo(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelInfo, Output: &buf})

	logr.Info("Info message")
	output := buf.String()
	if !strings.Contains(output, "I Info message") {
		t.Errorf("Expected info message in output, got %q", output)
	}

}

func TestInfof(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelInfo, Output: &buf})

	logr.Infof("Info message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "I Info message: formatted") {
		t.Errorf("Expected formatted info message in output, got %q", output)
	}

}

func TestInfoLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelInfo, Output: &buf})

	logr.Trace("trace message")
	logr.Debug("debug message")
	logr.Verbose("verbose message")
	if buf.Len() != 0 {
		t.Error("Nothing above Info should be logged at Info level")
	}

	logr.Info("info message")
	if buf.Len() == 0 {
		t.Error("Info should be logged at Info level")
	}

	buf.Reset()
	logr.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("Warn should be logged at Info level")
	}

	buf.Reset()
	logr.Error("error message")
	if buf.Len() == 0 {
		t.Error("Error should be logged at Info level")
	}

}

func TestTrace(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelTrace, LevelSet: true, Output: &buf})

	logr.Trace("Trace message")
	output := buf.String()
	if !strings.Contains(output, "T Trace message") {
		t.Errorf("Expected trace message in output, got %q", output)
	}

}

func TestTracef(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelTrace, LevelSet: true, Output: &buf})

	logr.Tracef("Trace message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "T Trace message: formatted") {
		t.Errorf("Expected formatted trace message in output, got %q", output)
	}

}

func TestTraceLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelTrace, LevelSet: true, Output: &buf})

	logr.Trace("trace message")
	if buf.Len() == 0 {
		t.Error("Trace should be logged at Trace level")
	}

	buf.Reset()
	logr.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("Debug should be logged at Trace level")
	}

	buf.Reset()
	logr.Info("info message")
	if buf.Len() == 0 {
		t.Error("Info should be logged at Trace level")
	}

}

func TestVerbose(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelVerbose, LevelSet: true, Output: &buf})

	logr.Verbose("Verbose message")
	output := buf.String()
	if !strings.Contains(output, "V Verbose message") {
		t.Errorf("Expected verbose message in output, got %q", output)
	}

}

func TestVerbosef(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelVerbose, LevelSet: true, Output: &buf})

	logr.Verbosef("Verbose message %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "V Verbose message formatted") {
		t.Errorf("Expected formatted verbose message in output, got %q", output)
	}

}

func TestVerboseLevel(t *testing.T) {

	// LevelCount 1: Verbose should log, Verbose2/3 should not
	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelVerbose, LevelSet: true, LevelCount: 1, Output: &buf})

	logr.Trace("trace message")
	logr.Debug("debug message")
	if buf.Len() != 0 {
		t.Error("Trace/Debug should not be logged at Verbose level")
	}

	logr.Verbose("verbose message")
	if buf.Len() == 0 {
		t.Error("Verbose should be logged at Verbose level")
	}

	buf.Reset()
	logr.Verbose2("verbose2 message")
	if buf.Len() != 0 {
		t.Error("Verbose2 should not be logged at LevelCount 1")
	}

	logr.Verbose3("verbose3 message")
	if buf.Len() != 0 {
		t.Error("Verbose3 should not be logged at LevelCount 1")
	}

	// LevelCount 2: Verbose and Verbose2 should log, Verbose3 should not
	buf.Reset()
	logr = NewWithOptions(context.TODO(), Options{Level: LogLevelVerbose, LevelSet: true, LevelCount: 2, Output: &buf})

	logr.Verbose("verbose message")
	if buf.Len() == 0 {
		t.Error("Verbose should be logged at LevelCount 2")
	}

	buf.Reset()
	logr.Verbose2("verbose2 message")
	if buf.Len() == 0 {
		t.Error("Verbose2 should be logged at LevelCount 2")
	}

	buf.Reset()
	logr.Verbose3("verbose3 message")
	if buf.Len() != 0 {
		t.Error("Verbose3 should not be logged at LevelCount 2")
	}

	// LevelCount 3: all Verbose levels should log
	buf.Reset()
	logr = NewWithOptions(context.TODO(), Options{Level: LogLevelVerbose, LevelSet: true, LevelCount: 3, Output: &buf})

	logr.Verbose3("verbose3 message")
	if buf.Len() == 0 {
		t.Error("Verbose3 should be logged at LevelCount 3")
	}

}

func TestDebugLevelCount(t *testing.T) {

	// LevelCount 1: Debug should log, Debug2/3 should not
	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, LevelCount: 1, Output: &buf})

	logr.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("Debug should be logged at LevelCount 1")
	}

	buf.Reset()
	logr.Debug2("debug2 message")
	if buf.Len() != 0 {
		t.Error("Debug2 should not be logged at LevelCount 1")
	}

	logr.Debug3("debug3 message")
	if buf.Len() != 0 {
		t.Error("Debug3 should not be logged at LevelCount 1")
	}

	// LevelCount 2: Debug and Debug2 should log
	buf.Reset()
	logr = NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, LevelCount: 2, Output: &buf})

	logr.Debug2("debug2 message")
	if buf.Len() == 0 {
		t.Error("Debug2 should be logged at LevelCount 2")
	}

	buf.Reset()
	logr.Debug3("debug3 message")
	if buf.Len() != 0 {
		t.Error("Debug3 should not be logged at LevelCount 2")
	}

	// LevelCount 3: all Debug levels should log
	buf.Reset()
	logr = NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, LevelCount: 3, Output: &buf})

	logr.Debug3("debug3 message")
	if buf.Len() == 0 {
		t.Error("Debug3 should be logged at LevelCount 3")
	}

}

func TestWarn(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelWarn, LevelSet: true, Output: &buf})

	logr.Warn("Warning message")
	output := buf.String()
	if !strings.Contains(output, "W Warning message") {
		t.Errorf("Expected warning message in output, got %q", output)
	}

}

func TestWarnf(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelWarn, LevelSet: true, Output: &buf})

	logr.Warnf("Warning message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "W Warning message: formatted") {
		t.Errorf("Expected formatted warning message in output, got %q", output)
	}

}

func TestWarnLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelWarn, LevelSet: true, Output: &buf})

	logr.Trace("trace message")
	logr.Debug("debug message")
	logr.Verbose("verbose message")
	logr.Info("info message")
	if buf.Len() != 0 {
		t.Error("Nothing below Warn should be logged at Warn level")
	}

	logr.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("Warn should be logged at Warn level")
	}

	buf.Reset()
	logr.Error("error message")
	if buf.Len() == 0 {
		t.Error("Error should be logged at Warn level")
	}

}

func TestJSON(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(context.TODO(), Options{Level: LogLevelDebug, LevelSet: true, Output: &buf})

	jsonString := `{"key": "value", "number": 123}`
	var jsonData map[string]interface{}

	json.Unmarshal([]byte(jsonString), &jsonData)

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")

	if err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}

	logr.Infof("Pretty JSON:\n%s", string(prettyJSON))

	output := buf.String()
	if !strings.Contains(output, "Pretty JSON:") {
		t.Errorf("Expected JSON output, got %q", output)
	}

}

func TestGetLoggerValuesFromString(t *testing.T) {

	tests := []struct {
		input      string
		wantLevel  LogLevel
		wantCount  int
	}{
		{"info", LogLevelInfo, 1},
		{"verbose", LogLevelVerbose, 1},
		{"verbose2", LogLevelVerbose, 2},
		{"verbose3", LogLevelVerbose, 3},
		{"debug", LogLevelDebug, 1},
		{"debug2", LogLevelDebug, 2},
		{"debug3", LogLevelDebug, 3},
		{"warn", LogLevelWarn, 1},
		{"warning", LogLevelWarn, 1},
		{"error", LogLevelError, 1},
		{"fatal", LogLevelFatal, 1},
		{"unknown", LogLevelInfo, 1},
		{"DEBUG", LogLevelDebug, 1},
	}

	for _, tt := range tests {
		level, count := GetLoggerValuesFromString(tt.input)
		if level != tt.wantLevel {
			t.Errorf("GetLoggerValuesFromString(%q): level = %v, want %v", tt.input, level, tt.wantLevel)
		}
		if count != tt.wantCount {
			t.Errorf("GetLoggerValuesFromString(%q): count = %d, want %d", tt.input, count, tt.wantCount)
		}
	}

}

func TestBuildMessageNonStringArgs(t *testing.T) {

	// Ensure buildMessage handles non-string arguments without panicking
	message := buildMessage(LogLevelInfo, 42, true, 3.14)
	if !strings.Contains(message, "42") {
		t.Errorf("Expected message to contain '42', got %q", message)
	}

}

func TestFromContextNil(t *testing.T) {

	ctx := context.TODO()
	logr := FromContext(ctx)
	if logr != nil {
		t.Error("Expected nil logger from empty context")
	}

}

func TestWithContextSameLogger(t *testing.T) {

	ctx := context.TODO()
	logr := New(ctx)
	ctx = WithContext(ctx, logr)

	// Calling WithContext with the same logger should return the same context
	ctx2 := WithContext(ctx, logr)
	if ctx2 != ctx {
		t.Error("Expected same context when logger already exists")
	}

}
