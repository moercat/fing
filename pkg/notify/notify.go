// Package notify 提供站内通知（in-app notification）。
//
// 设计：极简内存实现，重启清空。生产可改为 MySQL/Redis Stream。
//
// 用法：
//
//	notify.Push(userID, notify.Entry{
//	    Title: "新消息",
//	    Body:  "alice 给您发了私信",
//	    Type:  "message",
//	})
//
// 查询：
//
//	list := notify.List(userID, 20)  // 最近 20 条
//	unread := notify.UnreadCount(userID)
package notify

import (
	"sort"
	"sync"
	"time"
)

// Entry 一条通知
type Entry struct {
	ID       string    `json:"id"`
	UserID   uint      `json:"user_id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Type     string    `json:"type"`     // message / system / order ...
	Link     string    `json:"link"`     // 点击跳转的 URL
	Read     bool      `json:"read"`
	Time     time.Time `json:"time"`
}

var (
	mu      sync.RWMutex
	store   = make(map[uint][]Entry)
	counter int64
)

// Push 推送一条通知给指定用户
func Push(userID uint, e Entry) {
	mu.Lock()
	defer mu.Unlock()
	counter++
	e.ID = formatID(counter)
	e.UserID = userID
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	store[userID] = append(store[userID], e)
}

// List 返回用户最近 n 条通知（按时间倒序）
func List(userID uint, n int) []Entry {
	mu.RLock()
	defer mu.RUnlock()
	all := append([]Entry{}, store[userID]...)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Time.After(all[j].Time)
	})
	if n > 0 && n < len(all) {
		all = all[:n]
	}
	return all
}

// UnreadCount 未读数
func UnreadCount(userID uint) int {
	mu.RLock()
	defer mu.RUnlock()
	c := 0
	for _, e := range store[userID] {
		if !e.Read {
			c++
		}
	}
	return c
}

// MarkAllRead 全部标为已读
func MarkAllRead(userID uint) {
	mu.Lock()
	defer mu.Unlock()
	for i := range store[userID] {
		store[userID][i].Read = true
	}
}

func formatID(n int64) string {
	return time.Now().Format("20060102150405") + "-" + itoa(n)
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}