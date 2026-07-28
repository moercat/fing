// Package audit 提供操作日志（审计）能力。
//
// 默认是内存环形缓冲（最近 1000 条），可替换为：
//   - 文件（JSONL）
//   - MySQL
//   - Elasticsearch
//
// 用法：
//
//	audit.Log(audit.Entry{Method: "POST", Path: "/login", ...})
//
// 或通过 middleware.Audit() 自动记录。
package audit

import (
	"sync"
	"time"
)

// Entry 一条审计记录
type Entry struct {
	Time     time.Time     `json:"time"`
	Method   string        `json:"method"`
	Path     string        `json:"path"`
	IP       string        `json:"ip"`
	UserID   uint          `json:"user_id"`
	Status   int           `json:"status"`
	Latency  time.Duration `json:"latency_ms"`
	TraceID  string        `json:"trace_id"`
	UA       string        `json:"ua"`
	BodySize int           `json:"body_size"`
}

// Sink 是审计日志的存储接口
type Sink interface {
	Write(Entry)
}

// memSink 内存环形缓冲
type memSink struct {
	mu    sync.RWMutex
	buf   []Entry
	cap   int
	index int
}

func (m *memSink) Write(e Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.buf) < m.cap {
		m.buf = append(m.buf, e)
		return
	}
	m.buf[m.index] = e
	m.index = (m.index + 1) % m.cap
}

func (m *memSink) Snapshot() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, len(m.buf))
	copy(out, m.buf)
	return out
}

var defaultSink Sink = &memSink{cap: 1000}

// SetSink 替换默认 sink（如改为 MySQL sink）
func SetSink(s Sink) {
	if s != nil {
		defaultSink = s
	}
}

// Log 写入一条审计
func Log(e Entry) {
	defaultSink.Write(e)
}

// Recent 返回最近 n 条
func Recent(n int) []Entry {
	if ms, ok := defaultSink.(*memSink); ok {
		all := ms.Snapshot()
		if n > len(all) {
			n = len(all)
		}
		return all[len(all)-n:]
	}
	return nil
}