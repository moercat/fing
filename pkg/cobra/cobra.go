// Package cobra 提供轻量级定时任务注册与执行。
//
// 设计目标：
//   - 在 main.go 启动时调用 Cobra()，自动启动所有已注册任务
//   - 每个任务是一个 func() 闭包，无需继承任何接口
//   - 任务用 Register() 注册，集中管理
//   - panic 不影响其他任务（recover 包裹）
//
// 使用方式见 examples/task/ 目录。
package cobra

import (
	"fmt"
	"sync"
	"time"
)

// Task 一个定时任务的元信息
type Task struct {
	Name     string        // 任务名（日志用）
	Interval time.Duration // 间隔
	Fn       func()        // 任务函数
}

var (
	tasks []Task
	once  sync.Once
)

// Register 注册一个定时任务。
// 在 main() 里调用，Cobra() 启动时会按 Interval 周期执行。
//
// Example:
//
//	cobra.Register("hourly-cleanup", time.Hour, func() {
//	    // do something
//	})
func Register(name string, interval time.Duration, fn func()) {
	tasks = append(tasks, Task{
		Name:     name,
		Interval: interval,
		Fn:       fn,
	})
}

// Cobra 启动所有已注册任务。多次调用安全（只启动一次）。
// 通常在 main.go 里调一次。
func Cobra() {
	once.Do(func() {
		for _, t := range tasks {
			go run(t)
		}
		fmt.Printf("[cobra] started %d tasks\n", len(tasks))
	})
}

func run(t Task) {
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()
	for range ticker.C {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[cobra] task %q panicked: %v\n", t.Name, r)
				}
			}()
			t.Fn()
		}()
	}
}
