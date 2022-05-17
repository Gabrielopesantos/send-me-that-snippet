package logger

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

// Log levels
const (
	Info Level = iota
	Error
	Fatal
)

func (l Level) String() string {
	switch l {
	case Info:
		return "INFO"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{out: out, minLevel: minLevel}
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
		Trace      string            `json:"trace"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= Error {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(Error.String() + ": unable to marshal log message:" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// Might not be needed
//func (l *Logger) Write(message []byte) (n int, err error) {
//	return l.print(Error, string(message), nil)
//}

func (l *Logger) Info(message string, properties map[string]string) (n int, err error) {
	return l.print(Info, message, properties)
}

func (l *Logger) Error(message string, properties map[string]string) (n int, err error) {
	return l.print(Error, message, properties)
}

func (l *Logger) Fatal(message string, properties map[string]string) (n int, err error) {
	return l.print(Fatal, message, properties)
}
