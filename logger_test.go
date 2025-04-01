package logger

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code if needed
	// ...

	// Initialize logger
	logr := New()
	logr.SetLevel(LogLevelDebug)

	// Run tests
	exitCode := m.Run()

	// Teardown code if needed
	// ...

	os.Exit(exitCode)
}

func TestNewLogger(t *testing.T) {

	opts := Options{}

	logr := NewWithOptions(opts)

	if logr == nil {
		t.Fatal("Expected a valid Logger instance, got nil")
	}

	if logr.Level != LogLevelInfo {
		t.Errorf("Expected default log level to be LogLevelInfo, got %v", logr.Level)
	}

	if logr.Output != os.Stderr {
		t.Errorf("Expected default output to be os.Stderr, got %v", logr.Output)
	}

	return
}

func TestNewWithOptions(t *testing.T) {

	opts := Options{
		Level:  LogLevelDebug,
		Output: os.Stdout,
	}

	logr := NewWithOptions(opts)

	if logr.Level != LogLevelDebug {
		t.Errorf("Expected log level to be LogLevelDebug, got %v", logr.Level)
	}

	if logr.Output != os.Stdout {
		t.Errorf("Expected output to be os.Stdout, got %v", logr.Output)
	}

	return
}

func TestSetLevel(t *testing.T) {

	logr := New()
	logr.SetLevel(LogLevelError)

	if logr.Level != LogLevelError {
		t.Errorf("Expected log level to be LogLevelError, got %v", logr.Level)
	}

	return
}

func TestSetOutput(t *testing.T) {

	logr := New()
	tempFile, err := os.CreateTemp("", "testlog")

	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	logr.SetOutput(tempFile)

	if logr.Output != tempFile {
		t.Errorf("Expected output to be temp file, got %v", logr.Output)
	}

	return
}

func TestDebug(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelDebug, Output: os.Stderr})

	logr.Debug("Debug message")
	logr.Debug("Debug message\nwith newline")

	return
}

func TestDebugf(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelDebug, Output: os.Stderr})

	logr.Debugf("Debug message %s", "formatted")
	logr.Debugf("Debug message %s\nwith newline", "formatted")

	return
}

func TestDebugLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelDebug, Output: os.Stderr})

	logr.Trace("TestDebugLevel trace message")
	logr.Debug("TestDebugLevel debug message")
	logr.Verbose("TestDebugLevel verbose message")
	logr.Info("TestDebugLevel info message")
	logr.Warn("TestDebugLevel warn message")
	logr.Error("TestDebugLevel error message")

	return
}

func TestError(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelError, Output: os.Stderr})

	logr.Error("Error message")
	logr.Error("Error message\nwith newline")

	return
}

func TestErrorf(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelError, Output: os.Stderr})

	logr.Errorf("Error message: %s", "formatted")
	logr.Errorf("Error message: %s\nwith newline", "formatted")

	return
}

func TestErrorLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelError, Output: os.Stderr})

	logr.Trace("TestErrorLevel trace message")
	logr.Debug("TestErrorLevel debug message")
	logr.Verbose("TestErrorLevel verbose message")
	logr.Info("TestErrorLevel info message")
	logr.Warn("TestErrorLevel warn message")
	logr.Error("TestErrorLevel error message")

	return
}

func TestInfo(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelInfo, Output: os.Stderr})

	logr.Info("Info message")
	logr.Info("Info message\nwith newline")
	logr.Info("Info message\nwith newline\nwith more newlines")

	return
}

func TestInfof(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelInfo, Output: os.Stderr})

	logr.Infof("Info message: %s", "formatted")
	logr.Infof("Info message: %s\nwith newline", "formatted")

	return
}

func TestInfoLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelInfo, Output: os.Stderr})

	logr.Trace("TestInfoLevel trace message")
	logr.Debug("TestInfoLevel debug message")
	logr.Verbose("TestInfoLevel verbose message")
	logr.Info("TestInfoLevel info message")
	logr.Warn("TestInfoLevel warn message")
	logr.Error("TestInfoLevel error message")

	return
}

func TestTrace(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelTrace, Output: os.Stderr})

	logr.Trace("Trace message")
	logr.Trace("Trace message\nwith newline")

	return
}

func TestTracef(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelTrace, Output: os.Stderr})

	logr.Tracef("Trace message: %s", "formatted")
	logr.Tracef("Trace message: %s\nwith newline", "formatted")

	return
}

func TestTraceLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelTrace, Output: os.Stderr})

	logr.Trace("TestTraceLevel trace message")
	logr.Debug("TestTraceLevel debug message")
	logr.Verbose("TestTraceLevel verbose message")
	logr.Info("TestTraceLevel info message")
	logr.Warn("TestTraceLevel warn message")
	logr.Error("TestTraceLevel error message")

	return
}

func TestVerbose(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelVerbose, Output: os.Stderr})

	logr.Verbose("Verbose message")
	logr.Verbose("Verbose message\nwith newline")

	return
}

func TestVerbosef(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelVerbose, Output: os.Stderr})

	logr.Verbosef("Verbose message %s", "formatted")
	logr.Verbosef("Verbose message %s\nwith newline", "formatted")

	return
}

func TestVerboseLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelVerbose, Output: os.Stderr})

	logr.Trace("TestVerboseLevel trace message")
	logr.Debug("TestVerboseLevel debug message")
	logr.Verbose("TestVerboseLevel verbose message")
	logr.Info("TestVerboseLevel info message")
	logr.Warn("TestVerboseLevel warn message")
	logr.Error("TestVerboseLevel error message")

	return
}

func TestWarn(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelWarn, Output: os.Stderr})

	logr.Warn("Warning message")
	logr.Warn("Warning message\nwith newline")

	return
}

func TestWarnf(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelWarn, Output: os.Stderr})

	logr.Warnf("Warning message: %s", "formatted")
	logr.Warnf("Warning message: %s\nwith newline", "formatted")

	return
}

func TestWarnLevel(_ /*t*/ *testing.T) {

	logr := NewWithOptions(Options{Level: LogLevelWarn, Output: os.Stderr})

	logr.Trace("TestWarnLevel trace message")
	logr.Debug("TestWarnLevel debug message")
	logr.Verbose("TestWarnLevel verbose message")
	logr.Info("TestWarnLevel info message")
	logr.Warn("TestWarnLevel warn message")
	logr.Error("TestWarnLevel error message")

	return
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

	logr := NewWithOptions(Options{Level: LogLevelDebug, Output: os.Stderr})

	jsonString := "{\"key\": \"value\", \"number\": 123, \"boolean\": true, \"array\": [1, 2, 3], \"object\": {\"nestedKey\": \"nestedValue\"}}"
	var jsonData map[string]interface{}

	json.Unmarshal([]byte(jsonString), &jsonData)

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")

	if err != nil {
		t.Fatalf("Failed to format JSON: %v", err)
	}

	logr.Infof("Pretty JSON:\n%s", string(prettyJSON))

	return
}
