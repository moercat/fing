package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func init() {
	logDir := "./log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	fileName := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info 记录信息日志
func Info(v ...interface{}) {
	InfoLogger.Output(2, fmt.Sprint(v...))
}

// Infof 记录格式化信息日志
func Infof(format string, v ...interface{}) {
	InfoLogger.Output(2, fmt.Sprintf(format, v...))
}

// Warn 记录警告日志
func Warn(v ...interface{}) {
	WarnLogger.Output(2, fmt.Sprint(v...))
}

// Warnf 记录格式化警告日志
func Warnf(format string, v ...interface{}) {
	WarnLogger.Output(2, fmt.Sprintf(format, v...))
}

// Error 记录错误日志
func Error(v ...interface{}) {
	ErrorLogger.Output(2, fmt.Sprint(v...))
}

// Errorf 记录格式化错误日志
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Output(2, fmt.Sprintf(format, v...))
}

// Debug 记录调试日志
func Debug(v ...interface{}) {
	DebugLogger.Output(2, fmt.Sprint(v...))
}

// Debugf 记录格式化调试日志
func Debugf(format string, v ...interface{}) {
	DebugLogger.Output(2, fmt.Sprintf(format, v...))
}
