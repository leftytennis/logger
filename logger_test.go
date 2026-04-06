package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestNewLogger(t *testing.T) {

	opts := Options{}

	logr := NewWithOptions(opts)

	if logr == nil {
		t.Fatal("Expected a valid Logger instance, got nil")
	}

	ctx := WithContext(context.TODO(), logr)

	logrValue := FromContext(ctx)
	if logrValue == nil {
		t.Fatal("Expected a valid Logger instance from context, got nil")
	}

	if logr.GetLevel() != LevelInfo {
		t.Errorf("Expected default log level to be LevelInfo, got %v", logr.GetLevel())
	}

}

func TestNewWithOptions(t *testing.T) {

	opts := Options{
		Level:  LevelDebug,
		Output: os.Stdout,
	}

	logr := NewWithOptions(opts)

	if logr.GetLevel() != LevelDebug {
		t.Errorf("Expected log level to be LevelDebug, got %v", logr.GetLevel())
	}

}

func TestSetLevel(t *testing.T) {

	logr := New()
	var buf bytes.Buffer
	logr.SetOutput(&buf)
	logr.SetLevel(LevelError, 1)

	if logr.GetLevel() != LevelError {
		t.Errorf("Expected log level to be LevelError, got %v", logr.GetLevel())
	}

}

func TestSetLevelNoOutput(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	logr.SetLevel(LevelError, 1)

	if buf.Len() != 0 {
		t.Errorf("Expected no output from SetLevel, got %q", buf.String())
	}

}

func TestSetOutput(t *testing.T) {

	logr := New()
	var buf bytes.Buffer

	logr.SetOutput(&buf)

	logr.Info("test output")

	if buf.Len() == 0 {
		t.Error("Expected output to be written to buffer")
	}

}

func TestDebug(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelDebug, Output: &buf})

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
	logr := NewWithOptions(Options{Level: LevelDebug, Output: &buf})

	logr.Debugf("Debug message %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "D Debug message formatted") {
		t.Errorf("Expected formatted debug message in output, got %q", output)
	}

}

func TestDebugLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelDebug, LevelCount: 1, Output: &buf})

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
	logr := NewWithOptions(Options{Level: LevelError, Output: &buf})

	logr.Error("Error message")
	output := buf.String()
	if !strings.Contains(output, "E Error message") {
		t.Errorf("Expected error message in output, got %q", output)
	}

}

func TestErrorf(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelError, Output: &buf})

	logr.Errorf("Error message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "E Error message: formatted") {
		t.Errorf("Expected formatted error message in output, got %q", output)
	}

}

func TestErrorLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelError, Output: &buf})

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
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	logr.Info("Info message")
	output := buf.String()
	if !strings.Contains(output, "I Info message") {
		t.Errorf("Expected info message in output, got %q", output)
	}

}

func TestInfof(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	logr.Infof("Info message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "I Info message: formatted") {
		t.Errorf("Expected formatted info message in output, got %q", output)
	}

}

func TestInfoLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

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
	logr := NewWithOptions(Options{Level: LevelTrace, Output: &buf})

	logr.Trace("Trace message")
	output := buf.String()
	if !strings.Contains(output, "T Trace message") {
		t.Errorf("Expected trace message in output, got %q", output)
	}

}

func TestTracef(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, Output: &buf})

	logr.Tracef("Trace message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "T Trace message: formatted") {
		t.Errorf("Expected formatted trace message in output, got %q", output)
	}

}

func TestTraceLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, Output: &buf})

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

func TestTraceLevelCount(t *testing.T) {

	// LevelCount 1: Trace should log, Trace2/3 should not
	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, LevelCount: 1, Output: &buf})

	logr.Trace("trace message")
	if buf.Len() == 0 {
		t.Error("Trace should be logged at LevelCount 1")
	}

	buf.Reset()
	logr.Trace2("trace2 message")
	if buf.Len() != 0 {
		t.Error("Trace2 should not be logged at LevelCount 1")
	}

	logr.Trace3("trace3 message")
	if buf.Len() != 0 {
		t.Error("Trace3 should not be logged at LevelCount 1")
	}

	// LevelCount 2: Trace and Trace2 should log, Trace3 should not
	buf.Reset()
	logr = NewWithOptions(Options{Level: LevelTrace, LevelCount: 2, Output: &buf})

	logr.Trace2("trace2 message")
	if buf.Len() == 0 {
		t.Error("Trace2 should be logged at LevelCount 2")
	}

	buf.Reset()
	logr.Trace3("trace3 message")
	if buf.Len() != 0 {
		t.Error("Trace3 should not be logged at LevelCount 2")
	}

	// LevelCount 3: all Trace levels should log
	buf.Reset()
	logr = NewWithOptions(Options{Level: LevelTrace, LevelCount: 3, Output: &buf})

	logr.Trace3("trace3 message")
	if buf.Len() == 0 {
		t.Error("Trace3 should be logged at LevelCount 3")
	}

}

func TestTracef2(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, LevelCount: 2, Output: &buf})

	logr.Tracef2("trace2 %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "T trace2 formatted") {
		t.Errorf("Expected formatted trace2 message, got %q", output)
	}

}

func TestTracef3(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, LevelCount: 3, Output: &buf})

	logr.Tracef3("trace3 %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "T trace3 formatted") {
		t.Errorf("Expected formatted trace3 message, got %q", output)
	}

}

func TestVerbose(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelVerbose, Output: &buf})

	logr.Verbose("Verbose message")
	output := buf.String()
	if !strings.Contains(output, "V Verbose message") {
		t.Errorf("Expected verbose message in output, got %q", output)
	}

}

func TestVerbosef(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelVerbose, Output: &buf})

	logr.Verbosef("Verbose message %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "V Verbose message formatted") {
		t.Errorf("Expected formatted verbose message in output, got %q", output)
	}

}

func TestVerboseLevel(t *testing.T) {

	// LevelCount 1: Verbose should log, Verbose2/3 should not
	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelVerbose, LevelCount: 1, Output: &buf})

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
	logr = NewWithOptions(Options{Level: LevelVerbose, LevelCount: 2, Output: &buf})

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
	logr = NewWithOptions(Options{Level: LevelVerbose, LevelCount: 3, Output: &buf})

	logr.Verbose3("verbose3 message")
	if buf.Len() == 0 {
		t.Error("Verbose3 should be logged at LevelCount 3")
	}

}

func TestDebugLevelCount(t *testing.T) {

	// LevelCount 1: Debug should log, Debug2/3 should not
	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelDebug, LevelCount: 1, Output: &buf})

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
	logr = NewWithOptions(Options{Level: LevelDebug, LevelCount: 2, Output: &buf})

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
	logr = NewWithOptions(Options{Level: LevelDebug, LevelCount: 3, Output: &buf})

	logr.Debug3("debug3 message")
	if buf.Len() == 0 {
		t.Error("Debug3 should be logged at LevelCount 3")
	}

}

func TestWarn(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelWarn, Output: &buf})

	logr.Warn("Warning message")
	output := buf.String()
	if !strings.Contains(output, "W Warning message") {
		t.Errorf("Expected warning message in output, got %q", output)
	}

}

func TestWarnf(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelWarn, Output: &buf})

	logr.Warnf("Warning message: %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "W Warning message: formatted") {
		t.Errorf("Expected formatted warning message in output, got %q", output)
	}

}

func TestWarnLevel(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelWarn, Output: &buf})

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
	logr := NewWithOptions(Options{Level: LevelDebug, Output: &buf})

	jsonString := `{"key": "value", "number": 123}`
	var jsonData map[string]any

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

func TestParseLevel(t *testing.T) {

	tests := []struct {
		input     string
		wantLevel Level
		wantCount int
		wantErr   bool
	}{
		{"info", LevelInfo, 1, false},
		{"verbose", LevelVerbose, 1, false},
		{"verbose1", LevelVerbose, 1, false},
		{"verbose2", LevelVerbose, 2, false},
		{"verbose3", LevelVerbose, 3, false},
		{"debug", LevelDebug, 1, false},
		{"debug1", LevelDebug, 1, false},
		{"debug2", LevelDebug, 2, false},
		{"debug3", LevelDebug, 3, false},
		{"trace", LevelTrace, 1, false},
		{"trace1", LevelTrace, 1, false},
		{"trace2", LevelTrace, 2, false},
		{"trace3", LevelTrace, 3, false},
		{"warn", LevelWarn, 1, false},
		{"warning", LevelWarn, 1, false},
		{"error", LevelError, 1, false},
		{"fatal", LevelFatal, 1, false},
		{"DEBUG", LevelDebug, 1, false},
		{"TRACE", LevelTrace, 1, false},
		{"unknown", LevelInfo, 1, true},
		{"garbage", LevelInfo, 1, true},
	}

	for _, tt := range tests {
		level, count, err := ParseLevel(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseLevel(%q): err = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if level != tt.wantLevel {
			t.Errorf("ParseLevel(%q): level = %v, want %v", tt.input, level, tt.wantLevel)
		}
		if count != tt.wantCount {
			t.Errorf("ParseLevel(%q): count = %d, want %d", tt.input, count, tt.wantCount)
		}
	}

}

func TestBuildMessageNonStringArgs(t *testing.T) {

	// Ensure buildMessage handles non-string arguments without panicking
	message := buildMessage(LevelInfo, 42, true, 3.14)
	if !strings.Contains(message, "42") {
		t.Errorf("Expected message to contain '42', got %q", message)
	}

}

func TestBuildMessageMultipleArgs(t *testing.T) {

	message := buildMessage(LevelInfo, "hello", "world")
	if !strings.Contains(message, "hello world") {
		t.Errorf("Expected message to contain space-separated args 'hello world', got %q", message)
	}

}

func TestBuildMessageBlankLines(t *testing.T) {

	message := buildMessage(LevelInfo, "line1\n\nline3")
	lines := strings.Split(message, "\n")
	// Should have: "...I line1", "...  " (blank interior with padding), "...  line3", ""
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines, got %d: %q", len(lines), message)
	}
	if !strings.Contains(lines[0], "line1") {
		t.Errorf("Expected first line to contain 'line1', got %q", lines[0])
	}
	// The blank interior line should exist (with padding only)
	if strings.Contains(lines[1], "line") {
		t.Errorf("Expected second line to be blank (padding only), got %q", lines[1])
	}
	if !strings.Contains(lines[2], "line3") {
		t.Errorf("Expected third line to contain 'line3', got %q", lines[2])
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
	logr := New()
	ctx = WithContext(ctx, logr)

	// Calling WithContext with the same logger should return the same context
	ctx2 := WithContext(ctx, logr)
	if ctx2 != ctx {
		t.Error("Expected same context when logger already exists")
	}

}

func TestFatal(t *testing.T) {

	var buf bytes.Buffer
	var exitCode int
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})
	logr.exitFunc = func(code int) { exitCode = code }

	logr.Fatal("fatal error")

	output := buf.String()
	if !strings.Contains(output, "F fatal error") {
		t.Errorf("Expected fatal message in output, got %q", output)
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

}

func TestFatalf(t *testing.T) {

	var buf bytes.Buffer
	var exitCode int
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})
	logr.exitFunc = func(code int) { exitCode = code }

	logr.Fatalf("fatal error: %s", "details")

	output := buf.String()
	if !strings.Contains(output, "F fatal error: details") {
		t.Errorf("Expected fatal message in output, got %q", output)
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

}

func TestEnabled(t *testing.T) {

	logr := NewWithOptions(Options{Level: LevelDebug, LevelCount: 2, Output: os.Stderr})

	if !logr.Enabled(LevelDebug, 1) {
		t.Error("Expected Debug/1 to be enabled")
	}
	if !logr.Enabled(LevelDebug, 2) {
		t.Error("Expected Debug/2 to be enabled")
	}
	if logr.Enabled(LevelDebug, 3) {
		t.Error("Expected Debug/3 to NOT be enabled")
	}
	if logr.Enabled(LevelTrace, 1) {
		t.Error("Expected Trace to NOT be enabled at Debug level")
	}
	if !logr.Enabled(LevelInfo, 1) {
		t.Error("Expected Info to be enabled at Debug level")
	}

}

func TestConcurrentAccess(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, LevelCount: 3, Output: &buf})

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			logr.Info("concurrent", n)
			logr.Debug2("concurrent debug2", n)
			logr.Trace3("concurrent trace3", n)
			logr.SetLevel(LevelTrace, 3)
			_ = logr.GetLevel()
			_ = logr.GetLevelCount()
			_ = logr.Enabled(LevelDebug, 2)
		}(i)
	}
	wg.Wait()

}

func TestSetOutputNil(t *testing.T) {

	logr := New()
	logr.SetOutput(nil)

	// Should not panic; nil falls back to stderr
	logr.Info("fallback to stderr")

}

func TestWriteRawBytes(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	// Write implements io.Writer — raw bytes pass through with a trailing newline
	n, err := logr.Write([]byte("raw message"))
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n == 0 {
		t.Error("Expected non-zero bytes written")
	}
	output := buf.String()
	if !strings.Contains(output, "raw message") {
		t.Errorf("Expected raw message in output, got %q", output)
	}
	if output[len(output)-1] != '\n' {
		t.Error("Expected trailing newline")
	}

}

func TestWriteAlreadyNewline(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	logr.Write([]byte("message\n"))
	output := buf.String()
	if strings.HasSuffix(output, "\n\n") {
		t.Error("Write should not double the trailing newline")
	}

}

func TestBuildMessageNoArgs(t *testing.T) {

	message := buildMessage(LevelInfo)
	if !strings.Contains(message, " I ") {
		t.Errorf("Expected Info prefix in empty message, got %q", message)
	}
	if message[len(message)-1] != '\n' {
		t.Error("Expected trailing newline")
	}

}

func TestLevelStringUnknown(t *testing.T) {

	unknown := Level(99)
	if unknown.String() != "Unknown" {
		t.Errorf("Expected 'Unknown' for invalid level, got %q", unknown.String())
	}

}

func TestParseLevelWhitespace(t *testing.T) {

	level, count, err := ParseLevel("  debug  ")
	if err != nil {
		t.Fatalf("ParseLevel with whitespace returned error: %v", err)
	}
	if level != LevelDebug {
		t.Errorf("Expected LevelDebug, got %v", level)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

}

func TestHigherLevelEnablesLowerSubLevels(t *testing.T) {

	// When set to Trace, all Debug sub-levels should be enabled
	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelTrace, LevelCount: 1, Output: &buf})

	logr.Debug2("debug2 at trace level")
	if buf.Len() == 0 {
		t.Error("Debug2 should be logged when logger is set to Trace (higher level)")
	}

	buf.Reset()
	logr.Debug3("debug3 at trace level")
	if buf.Len() == 0 {
		t.Error("Debug3 should be logged when logger is set to Trace (higher level)")
	}

	buf.Reset()
	logr.Verbose3("verbose3 at trace level")
	if buf.Len() == 0 {
		t.Error("Verbose3 should be logged when logger is set to Trace (higher level)")
	}

}

func TestWriteEmptyBytes(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	n, err := logr.Write([]byte{})
	if err != nil {
		t.Fatalf("Write of empty bytes returned error: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected 0 bytes written for empty input, got %d", n)
	}

}

func TestDebugf2(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelDebug, LevelCount: 2, Output: &buf})

	logr.Debugf2("debug2 %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "D debug2 formatted") {
		t.Errorf("Expected formatted debug2 message, got %q", output)
	}

}

func TestDebugf3(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelDebug, LevelCount: 3, Output: &buf})

	logr.Debugf3("debug3 %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "D debug3 formatted") {
		t.Errorf("Expected formatted debug3 message, got %q", output)
	}

}

func TestVerbosef2(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelVerbose, LevelCount: 2, Output: &buf})

	logr.Verbosef2("verbose2 %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "V verbose2 formatted") {
		t.Errorf("Expected formatted verbose2 message, got %q", output)
	}

}

func TestVerbosef3(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelVerbose, LevelCount: 3, Output: &buf})

	logr.Verbosef3("verbose3 %s", "formatted")
	output := buf.String()
	if !strings.Contains(output, "V verbose3 formatted") {
		t.Errorf("Expected formatted verbose3 message, got %q", output)
	}

}

func TestFatalLevel(t *testing.T) {

	// Fatal should always log regardless of the configured level
	var buf bytes.Buffer
	var exitCode int
	logr := NewWithOptions(Options{Level: LevelFatal, Output: &buf})
	logr.exitFunc = func(code int) { exitCode = code }

	logr.Fatal("fatal at fatal level")

	if buf.Len() == 0 {
		t.Error("Fatal should always be logged")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

}

func TestLogMultipleArgs(t *testing.T) {

	var buf bytes.Buffer
	logr := NewWithOptions(Options{Level: LevelInfo, Output: &buf})

	logr.Info("count:", 42, "flag:", true)
	output := buf.String()
	if !strings.Contains(output, "count: 42 flag: true") {
		t.Errorf("Expected space-separated args, got %q", output)
	}

}

func TestGetLevelCount(t *testing.T) {

	logr := NewWithOptions(Options{Level: LevelDebug, LevelCount: 2})
	if logr.GetLevelCount() != 2 {
		t.Errorf("Expected level count 2, got %d", logr.GetLevelCount())
	}

	logr.SetLevel(LevelTrace, 3)
	if logr.GetLevelCount() != 3 {
		t.Errorf("Expected level count 3 after SetLevel, got %d", logr.GetLevelCount())
	}

}

func TestParseLevelAllLevels(t *testing.T) {

	// Verify all Level.String() round-trips through ParseLevel
	levels := []struct {
		level Level
		name  string
	}{
		{LevelFatal, "fatal"},
		{LevelError, "error"},
		{LevelWarn, "warn"},
		{LevelInfo, "info"},
		{LevelVerbose, "verbose"},
		{LevelDebug, "debug"},
		{LevelTrace, "trace"},
	}

	for _, tt := range levels {
		level, _, err := ParseLevel(tt.name)
		if err != nil {
			t.Errorf("ParseLevel(%q) returned error: %v", tt.name, err)
		}
		if level != tt.level {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.name, level, tt.level)
		}
		if level.String() != fmt.Sprintf("%s%s", strings.ToUpper(tt.name[:1]), tt.name[1:]) {
			t.Errorf("Level(%d).String() = %q, want title-cased %q", level, level.String(), tt.name)
		}
	}

}
