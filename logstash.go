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

var writeCount uint

// Fields ..
type Fields map[string]interface{}

// inst ..
type inst struct {
	fields Fields
	msg    string
	time   string
	level  Level
}

// Entry ..
type Entry interface {
	Panic(string)
	Fatal(string)
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
	WithFields(Fields) *inst
}

// WithFields add more field
func WithFields(fields Fields) (entry Entry) {
	entry = &inst{fields: fields}
	entry.WithFields(fields)
	return
}

// WithFields ..
func (i *inst) WithFields(fields Fields) *inst {
	for k, v := range i.fields {
		i.fields[k] = v
	}
	return i
}

// Panic ..
func Panic(str string) {
	i := &inst{}
	i.Panic(str)
}

// Panic ..
func (i *inst) Panic(str string) {
	i.msg = str
	i.level = PanicLevel
	i.output()
}

// Fatal ..
func Fatal(str string) {
	i := &inst{}
	i.Fatal(str)
}

// Fatal ..
func (i *inst) Fatal(str string) {
	i.msg = str
	i.level = FatalLevel
	i.output()
}

// Error ..
func Error(str string) {
	i := &inst{}
	i.Error(str)
}

// Error ..
func (i *inst) Error(str string) {
	i.msg = str
	i.level = ErrorLevel
	i.output()
}

// Warn ..
func Warn(str string) {
	i := &inst{}
	i.Warn(str)
}

// Warn ..
func (i *inst) Warn(str string) {
	i.msg = str
	i.level = WarnLevel
	i.output()
}

// Info ..
func Info(str string) {
	i := &inst{}
	i.Info(str)
}

// Info ..
func (i *inst) Info(str string) {
	i.msg = str
	i.level = InfoLevel
	i.output()
}

// Debug ..
func Debug(str string) {
	i := &inst{}
	i.Debug(str)
}

// Debug ..
func (i *inst) Debug(str string) {
	i.msg = str
	i.level = DebugLevel
	i.output()
}

func (i *inst) output() {
	var color int
	var err error
	var waitWrite []byte
	if i.level > config.LoggerLevel {
		return
	}
	switch i.level {
	case DebugLevel:
		color = gray
	case WarnLevel:
		color = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		color = red
	default:
		color = blue
	}
	levelText := strings.ToUpper(i.level.String())[0:4]
	var keys string
	if i.fields == nil {
		i.fields = Fields{}
	}
	if _, file, line, ok := runtime.Caller(2); ok {
		if filepath.Base(file) == "logstash.go" {
			if _, file, line, ok := runtime.Caller(3); ok {
				i.fields["__file"] = filepath.Base(file)
				i.fields["__line"] = line
			}
		} else {
			i.fields["__file"] = filepath.Base(file)
			i.fields["__line"] = line
		}
	}

	for k, v := range i.fields {
		if k != "__line" && k != "__file" {
			color := green
			keys += fmt.Sprintf(" \x1b[%dm%s\x1b[0m=%+v", color, k, v)
		}
	}

	for k, v := range i.fields {
		color := red
		if k == "__line" || k == "__file" {
			keys += fmt.Sprintf(" \x1b[%dm%s\x1b[0m=%+v", color, k[2:], v)
		}
	}

	i.time = time.Now().Format("Jan 2 15:04:05")
	if config.ConsoleOutput {
		fmt.Printf("\x1b[%dm%s\x1b[0m[%s] %-40s %s\n", color, levelText, i.time, i.msg, keys)
	}
	i.fields["time"] = i.time
	i.fields["level"] = i.level.String()
	i.fields["msg"] = i.msg
	if waitWrite, err = json.Marshal(i.fields); err != nil {
		Fatal("Cannot convert fields to string")
		return
	}
	waitWrite = append(waitWrite, '\n')

	if fileOutput {
		lock.Lock()
		if writeCount == 100 {
			resetLogFile()
			writeCount = 0
		}
		writeCount++
		logFile.Write(waitWrite)
		lock.Unlock()
	}

	if PanicLevel == i.level {
		os.Exit(-1)
	}
}
