package logger

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code if needed
	// ...

	// Initialize logger
	logr := NewWithOptions(ctx, Options{
		Level:  LogLevelInfo,
		LevelCount: 3,
		Output: os.Stderr,
	})
	
	ctx = WithContext(context.TODO(), logr)

	value := FromContext(ctx)

	logr.SetLevel(LogLevelDebug)
	logr.Verbosef2("Logger initialized in TestMain: %v", value)

	// Run tests
	exitCode := m.Run()

	// Teardown code if needed
	// ...

	os.Exit(exitCode)
}

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

	opts := Options{
		Level:  LogLevelDebug,
		Output: os.Stdout,
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
	logr.SetLevel(LogLevelError)

	if logr.Level != LogLevelError {
		t.Errorf("Expected log level to be LogLevelError, got %v", logr.Level)
	}

}

func TestSetOutput(t *testing.T) {

	logr := New(context.Background())
	tempFile, err := os.CreateTemp("", "testlog")

	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	logr.SetOutput(tempFile)

	if logr.Output != tempFile {
		t.Errorf("Expected output to be temp file, got %v", logr.Output)
	}

}

func TestDebug(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelDebug, Output: os.Stderr})

	logr.Debug("Debug message")
	logr.Debug("Debug message\nwith newline")
	logr.Debug("Debug message\nwith newline\nwith more newlines")

}

func TestDebugf(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelDebug, Output: os.Stderr})

	logr.Debugf("Debug message %s", "formatted")
	logr.Debugf("Debug message %s\nwith newline", "formatted")

}

func TestDebugLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelDebug, LevelCount: 1,Output: os.Stderr})

	logr.Trace("TestDebugLevel trace message")
	logr.Debug("TestDebugLevel debug message")
	logr.Verbose("TestDebugLevel verbose message")
	logr.Info("TestDebugLevel info message")
	logr.Warn("TestDebugLevel warn message")
	logr.Error("TestDebugLevel error message")

}

func TestError(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelError, Output: os.Stderr})

	logr.Error("Error message")
	logr.Error("Error message\nwith newline")

}

func TestErrorf(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelError, Output: os.Stderr})

	logr.Errorf("Error message: %s", "formatted")
	logr.Errorf("Error message: %s\nwith newline", "formatted")

}

func TestErrorLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelError, Output: os.Stderr})

	logr.Trace("TestErrorLevel trace message")
	logr.Debug("TestErrorLevel debug message")
	logr.Verbose("TestErrorLevel verbose message")
	logr.Info("TestErrorLevel info message")
	logr.Warn("TestErrorLevel warn message")
	logr.Error("TestErrorLevel error message")

}

func TestInfo(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelInfo, Output: os.Stderr})

	logr.Info("Info message")
	logr.Info("Info message\nwith newline")
	logr.Info("Info message\nwith newline\nwith more newlines")

}

func TestInfof(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelInfo, Output: os.Stderr})

	logr.Infof("Info message: %s", "formatted")
	logr.Infof("Info message: %s\nwith newline", "formatted")

}

func TestInfoLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelInfo, Output: os.Stderr})

	logr.Trace("TestInfoLevel trace message")
	logr.Debug("TestInfoLevel debug message")
	logr.Verbose("TestInfoLevel verbose message")
	logr.Info("TestInfoLevel info message")
	logr.Warn("TestInfoLevel warn message")
	logr.Error("TestInfoLevel error message")

}

func TestTrace(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelTrace, Output: os.Stderr})

	logr.Trace("Trace message")
	logr.Trace("Trace message\nwith newline")

}

func TestTracef(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelTrace, Output: os.Stderr})

	logr.Tracef("Trace message: %s", "formatted")
	logr.Tracef("Trace message: %s\nwith newline", "formatted")

}

func TestTraceLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelTrace, Output: os.Stderr})

	logr.Trace("TestTraceLevel trace message")
	logr.Debug("TestTraceLevel debug message")
	logr.Verbose("TestTraceLevel verbose message")
	logr.Info("TestTraceLevel info message")
	logr.Warn("TestTraceLevel warn message")
	logr.Error("TestTraceLevel error message")

}

func TestVerbose(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelVerbose, Output: os.Stderr})

	logr.Verbose("Verbose message")
	logr.Verbose("Verbose message\nwith newline")
	logr.Verbose("Verbose message2")
	logr.Verbose("Verbose message2\nwith newline")

}

func TestVerbosef(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelVerbose, Output: os.Stderr})

	logr.Verbosef("Verbose message %s", "formatted")
	logr.Verbosef("Verbose message %s\nwith newline", "formatted")

}

func TestVerboseLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelVerbose, Output: os.Stderr})

	logr.Trace("TestVerboseLevel trace message")
	logr.Debug("TestVerboseLevel debug message")
	logr.Verbose("TestVerboseLevel verbose message")
	logr.Info("TestVerboseLevel info message")
	logr.Warn("TestVerboseLevel warn message")
	logr.Error("TestVerboseLevel error message")

}

func TestWarn(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelWarn, Output: os.Stderr})

	logr.Warn("Warning message")
	logr.Warn("Warning message\nwith newline")

}

func TestWarnf(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelWarn, Output: os.Stderr})

	logr.Warnf("Warning message: %s", "formatted")
	logr.Warnf("Warning message: %s\nwith newline", "formatted")

}

func TestWarnLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelWarn, Output: os.Stderr})

	logr.Trace("TestWarnLevel trace message")
	logr.Debug("TestWarnLevel debug message")
	logr.Verbose("TestWarnLevel verbose message")
	logr.Info("TestWarnLevel info message")
	logr.Warn("TestWarnLevel warn message")
	logr.Error("TestWarnLevel error message")

}

// func TestFatal(_ /*t*/ *testing.T) {

// 	logr := NewWithOptions(Options{Level: LogLevelFatal, Output: os.Stderr})

// 	origFatal := logFatal
// 	defer func() { logFatal = origFatal }() // Restore original function

// 	logFatal = func(writer Logger, a ...any) {
// 		message := buildMessage(LogLevelFatal, a...)
// 		_, err := writer.Write([]byte(message))

// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	logr.Fatal("Fatal message")
// 	logr.Fatal("Fatal message\nwith newline")

// 	return
// }

func TestJSON(t *testing.T) {

	logr := NewWithOptions(ctx, Options{Level: LogLevelDebug, Output: os.Stderr})

	jsonString := "{\"key\": \"value\", \"number\": 123, \"boolean\": true, \"array\": [1, 2, 3], \"object\": {\"nestedKey\": \"nestedValue\"}}"
	var jsonData map[string]interface{}

	json.Unmarshal([]byte(jsonString), &jsonData)

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")

	if err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}

	logr.Infof("Pretty JSON:\n%s", string(prettyJSON))

}
