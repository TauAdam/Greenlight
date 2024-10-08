package json_logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level represents the severity level for a log entry.
type Level int8

const (
	LevelInfo  Level = iota // Has the value 0.
	LevelError              // Has the value 1.
	LevelFatal              // Has the value 2.
	LevelOff                // Has the value 3.
)

// String returns a human-friendly string for the log level
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Logger represents a custom JSON logger, which writes log entries at or above a minimum log
// level to a specific output destination.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	_, err := l.print(LevelInfo, message, properties)
	if err != nil {
		_, _ = l.print(LevelError, "failed to write log entry: "+err.Error(), nil)
	}
}
func (l *Logger) PrintError(inputErr error, properties map[string]string) {
	_, err := l.print(LevelError, inputErr.Error(), properties)
	if err != nil {
		_, _ = l.print(LevelError, "failed to write log entry: "+err.Error(), nil)
	}
}
func (l *Logger) PrintFatal(inputErr error, properties map[string]string) {
	_, err := l.print(LevelFatal, inputErr.Error(), properties)
	if err != nil {
		_, _ = l.print(LevelError, "failed to write log entry: "+err.Error(), nil)
	}
	// For entries at the FATAL level, terminate the application.
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}
	// Lock the mutex so that no two writes to the output destination cannot happen concurrently.
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// Write implements the io.Writer interface
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
