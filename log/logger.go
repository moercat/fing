package log

import (
	"fing/pkg/config"
	"fmt"
	"os"
	"time"
)

const (
	LevelError         = iota // LevelError 错误
	LevelWarning              // LevelWarning 警告
	LevelInformational        // LevelInformational 提示
	LevelDebug                // LevelDebug 除错
)

var logger = &Logger{
	level: config.Config.Level,
}

// Logger 日志
type Logger struct {
	level int
}

// Println 打印
func (ll *Logger) Println(msg string) {
	fmt.Printf("%s %s", time.Now().In(time.Local).Format("2006-01-02 15:04:05"), msg)
}

// Panic 极端错误
func Panic(format string, v ...interface{}) {
	if LevelError > logger.level {
		return
	}
	msg := fmt.Sprintf("[Panic] "+format+" \n", v...)
	logger.Println(msg)
	os.Exit(0)
}

// Error 错误
func Error(format string, v ...interface{}) {
	if LevelError > logger.level {
		return
	}
	msg := fmt.Sprintf("[E] "+format+" \n", v...)
	logger.Println(msg)
}

// Warning 警告
func Warning(format string, v ...interface{}) {
	if LevelWarning > logger.level {
		return
	}
	msg := fmt.Sprintf("[W] "+format+" \n", v...)
	logger.Println(msg)
}

// Info 信息
func Info(format string, v ...interface{}) {
	if LevelInformational > logger.level {
		return
	}
	msg := fmt.Sprintf("[I] "+format+" \n", v...)
	logger.Println(msg)
}

// Debug 校验
func Debug(format string, v ...interface{}) {
	if LevelDebug > logger.level {
		return
	}
	msg := fmt.Sprintf("[D] "+format+" \n", v...)
	logger.Println(msg)
}
