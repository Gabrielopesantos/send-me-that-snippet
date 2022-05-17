package logger

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type logMessageContent struct {
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace"`
}

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

type outputEncoding int8

const (
	Json outputEncoding = iota
	Console
)

func encodeJson(data *logMessageContent) []byte {
	var line []byte
	line, err := json.Marshal(data)
	if err != nil {
		line = []byte(Error.String() + ": unable to marshal log message:" + err.Error())
	}

	return line
}

func encodeConsole(data *logMessageContent) []byte {
	// In progress
	var line string
	line += data.Level + " | " + data.Time + " | " + data.Message

	return []byte(line)
}

var encodingFuncsMap = map[outputEncoding]func(*logMessageContent) []byte{
	Json:    encodeJson,
	Console: encodeConsole,
}

type Logger struct {
	out        io.Writer
	minLevel   Level
	encodeFunc func(data *logMessageContent) []byte
	mu         sync.Mutex
}

func New(out io.Writer, minLevel Level, encoding outputEncoding) *Logger {
	return &Logger{
		out:        out,
		minLevel:   minLevel,
		encodeFunc: encodingFuncsMap[encoding],
	}
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	logMsg := logMessageContent{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= Error {
		logMsg.Trace = string(debug.Stack())
	}

	log := l.encodeFunc(&logMsg)

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(log, '\n'))
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
