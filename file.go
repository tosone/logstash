package logstash

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

var fileOutput = false               // 是否可以输出到文件
var logFile *os.File                 // 日志文件的文件对象
var logFileName string               // 日志文件的文件名
var maxSize int64 = 10 * 1024 * 1024 // 文件的最大文件大小

var dir string

var prefix string

var lock sync.Mutex

// SetMaxSize 设置最大文件大小
func SetMaxSize(m int64) {
	maxSize = m * 1024 * 1024
}

// SetOutput 设置文件输出
func SetOutput(d, p string) {
	var err error
	dir = d
	prefix = p
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		WithFields(Fields{"dir": dir}).Fatal("Directory is not exist.")
		return
	}
	if err != nil {
		Fatal(err.Error())
		return
	}
	logFile, err = validLogFile(dir, path.Join(dir, prefix+".log"))
	if err != nil {
		Fatal(err.Error())
		return
	}
	fileOutput = true
}

func resetLogFile() {
	var err error
	lock.Lock()
	fileOutput = false
	if ok, err := checkFileSize(logFileName); ok && err == nil {
		return
	}
	logFile, err = validLogFile(dir, path.Join(dir, prefix+".log"))
	if err != nil {
		Fatal(err.Error())
		return
	}
	fileOutput = true
	lock.Unlock()
}

func validLogFile(dir, file string) (*os.File, error) {
	var fileInfo os.FileInfo
	var err error
	if fileInfo, err = os.Stat(file); os.IsNotExist(err) {
		return os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	}
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() > maxSize {
		if filepath.Ext(file) == ".log" {
			return validLogFile(dir, file+".1")
		}
		index, err := strconv.Atoi(filepath.Ext(file)[1:])
		if err != nil {
			return nil, err
		}
		return validLogFile(dir, filepath.Base(file)[0:len(filepath.Base(file))-len(filepath.Ext(file))]+"."+strconv.Itoa(index+1))
	}
	logFileName = file
	return os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
}

func checkFileSize(file string) (bool, error) {
	var fileInfo os.FileInfo
	var err error
	if file == "" {
		return false, nil
	}
	if fileInfo, err = os.Stat(file); err != nil {
		return false, err
	}
	if fileInfo.Size() > maxSize {
		return false, nil
	}
	return true, nil
}
