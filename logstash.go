package logstash

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

// Fields ..
type Fields map[string]interface{}

// Entry ..
type Entry struct {
	Fields Fields
	Msg    string
	Time   string
	Level  Level
}

// WithFields ..
func WithFields(fields Fields) *Entry {
	entry := &Entry{Fields: fields}
	entry.WithFields(fields)
	return entry
}

// WithFields ..
func (entry *Entry) WithFields(fields Fields) *Entry {
	for k, v := range entry.Fields {
		entry.Fields[k] = v
	}
	return entry
}

// Panic ..
func Panic(str string) {
	entry := &Entry{}
	entry.Panic(str)
}

// Panic ..
func (entry *Entry) Panic(str string) {
	entry.Msg = str
	entry.Level = PanicLevel
	entry.output()
}

// Fatal ..
func Fatal(str string) {
	entry := &Entry{}
	entry.Fatal(str)
}

// Fatal ..
func (entry *Entry) Fatal(str string) {
	entry.Msg = str
	entry.Level = FatalLevel
	entry.output()
}

// Error ..
func Error(str string) {
	entry := &Entry{}
	entry.Error(str)
}

// Error ..
func (entry *Entry) Error(str string) {
	entry.Msg = str
	entry.Level = ErrorLevel
	entry.output()
}

// Warn ..
func Warn(str string) {
	entry := &Entry{}
	entry.Warn(str)
}

// Warn ..
func (entry *Entry) Warn(str string) {
	entry.Msg = str
	entry.Level = WarnLevel
	entry.output()
}

// Info ..
func Info(str string) {
	entry := &Entry{}
	entry.Info(str)
}

// Info ..
func (entry *Entry) Info(str string) {
	entry.Msg = str
	entry.Level = InfoLevel
	entry.output()
}

// Debug ..
func Debug(str string) {
	entry := &Entry{}
	entry.Debug(str)
}

// Debug ..
func (entry *Entry) Debug(str string) {
	entry.Msg = str
	entry.Level = DebugLevel
	entry.output()
}

func (entry *Entry) output() {
	var color int
	var err error
	var waitWrite []byte
	switch entry.Level {
	case DebugLevel:
		color = gray
	case WarnLevel:
		color = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		color = red
	default:
		color = blue
	}
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	var keys string
	if entry.Fields == nil {
		entry.Fields = Fields{}
	}
	if _, file, line, ok := runtime.Caller(2); ok {
		if filepath.Base(file) == "logstash.go" {
			if _, file, line, ok := runtime.Caller(3); ok {
				entry.Fields["__file"] = filepath.Base(file)
				entry.Fields["__line"] = line
			}
		} else {
			entry.Fields["__file"] = filepath.Base(file)
			entry.Fields["__line"] = line
		}
	}
	for k, v := range entry.Fields {
		color := red
		if k != "__line" && k != "__file" {
			color = green
			keys += fmt.Sprintf(" \x1b[%dm%s\x1b[0m=%+v", color, k, v)
		} else {
			keys += fmt.Sprintf(" \x1b[%dm%s\x1b[0m=%+v", color, k[2:], v)
		}
	}
	entry.Time = time.Now().Format("Jan 2 15:04:05")
	fmt.Printf("\x1b[%dm%s\x1b[0m[%s] %-40s %s\n", color, levelText, entry.Time, entry.Msg, keys)
	entry.Fields["time"] = entry.Time
	entry.Fields["level"] = entry.Level.String()
	entry.Fields["msg"] = entry.Msg
	if waitWrite, err = json.Marshal(entry.Fields); err != nil {
		Fatal("Cannot convert fields to string")
		return
	}
	waitWrite = append(waitWrite, '\n')
	lock.Lock()
	logFile.Write(waitWrite)
	lock.Unlock()
	if PanicLevel == entry.Level {
		os.Exit(-1)
	}
}
