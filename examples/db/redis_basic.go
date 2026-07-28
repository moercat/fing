//go:build ignore
// +build ignore

// Package main 演示 fing 的 Redis（go-redis v7）使用方式。
//
// fing/pkg/db 提供了开箱即用的 RedisClient，业务代码直接用。
//
// fing 用的是 go-redis v7，API 形如：
//   db.RedisClient.Set(key, value, ttl).Err()
//   db.RedisClient.Get(key).Result()
//
// 实际想用 context 的话用 db.RedisClient.WithContext(ctx).
package main

import (
	"fmt"
	"time"

	"fing/pkg/db"
)

// ============================================================
// 字符串操作（最常用）
// ============================================================

// exampleString 字符串 SET/GET — 缓存场景
func exampleString() {
	key := "user:1:profile"

	// SET key value EX 60 (60 秒过期)
	if err := db.RedisClient.Set(key, "alice", 60*time.Second).Err(); err != nil {
		fmt.Println("set error:", err)
	}

	// GET key
	val, err := db.RedisClient.Get(key).Result()
	if err != nil {
		fmt.Println("get error:", err)
		return
	}
	fmt.Println("cached:", val)
}

// ============================================================
// Hash 操作（适合存对象）
// ============================================================

// exampleHash Hash — 存用户对象字段
func exampleHash() {
	key := "user:1"

	// 写字段
	if err := db.RedisClient.HSet(key, map[string]interface{}{
		"user_name": "alice",
		"nickname":  "Alice",
		"avatar":    "https://...",
	}).Err(); err != nil {
		fmt.Println("hset error:", err)
	}

	// 读单个字段
	nickname, _ := db.RedisClient.HGet(key, "nickname").Result()
	fmt.Println("nickname:", nickname)

	// 读所有字段
	all, _ := db.RedisClient.HGetAll(key).Result()
	fmt.Printf("user obj: %+v\n", all)
}

// ============================================================
// 计数器（INCR / DECR）
// ============================================================

// exampleCounter 计数器 — 限流、点赞、阅读量
func exampleCounter() {
	key := "article:123:views"

	// 自增 1，返回新值
	views, _ := db.RedisClient.Incr(key).Result()
	fmt.Println("views:", views)

	// 自增指定值
	db.RedisClient.IncrBy(key, 10)

	// 设置过期（24h 后清零重新计数）
	db.RedisClient.Expire(key, 24*time.Hour)
}

// ============================================================
// 列表（队列 / 最新列表）
// ============================================================

// exampleList 列表 — 最近消息队列
func exampleList() {
	key := "chat:room:1:messages"

	// LPUSH 左侧推入
	db.RedisClient.LPush(key, "msg-1")
	db.RedisClient.LPush(key, "msg-2", "msg-3")

	// LRANGE 取最近 10 条
	msgs, _ := db.RedisClient.LRange(key, 0, 9).Result()
	fmt.Println("recent:", msgs)

	// LTRIM 保留最新 100 条
	db.RedisClient.LTrim(key, 0, 99)
}

// ============================================================
// Set（去重 / 标签）
// ============================================================

func exampleSet() {
	key := "user:1:followers"

	db.RedisClient.SAdd(key, "alice", "bob", "carol")
	count, _ := db.RedisClient.SCard(key).Result()
	fmt.Println("follower count:", count)

	// 交集
	db.RedisClient.SInter("tag:go", "tag:backend")
}

// ============================================================
// 有序集合（排行榜）
// ============================================================

// exampleLeaderboard 排行榜
func exampleLeaderboard() {
	key := "leaderboard:2026-07"

	// 添加分数
	db.RedisClient.ZAdd(key, redisZ(100, "alice"), redisZ(95, "bob"))

	// 取前 10 名（带分数）
	top, _ := db.RedisClient.ZRevRangeWithScores(key, 0, 9).Result()
	for _, z := range top {
		fmt.Printf("%s: %v\n", z.Member, z.Score)
	}
}

func redisZ(score float64, member string) redisZItem {
	return redisZItem{Score: float64(score), Member: member}
}

type redisZItem struct {
	Score  float64
	Member interface{}
}

// ============================================================
// 分布式锁（单实例）
// ============================================================

// exampleLock 用 SETNX 实现分布式锁
func exampleLock() error {
	lockKey := "lock:order:123"

	// SETNX EX 5 (5 秒过期，防止死锁)
	ok, err := db.RedisClient.SetNX(lockKey, "owner-uuid", 5*time.Second).Result()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("lock held by others")
	}
	// 拿到锁，做事
	fmt.Println("got lock, doing work...")
	// 释放锁（用 lua 脚本保证原子性更安全，但简单场景直接 DEL 也行）
	defer db.RedisClient.Del(lockKey)
	return nil
}

// ============================================================
// 缓存模式（Cache-Aside）
// ============================================================

type userProfile struct {
	UserName string
	Nickname string
	Avatar   string
}

// getUserWithCache Cache-Aside 模式
func getUserWithCache(userID uint) (*userProfile, error) {
	cacheKey := fmt.Sprintf("user:profile:%d", userID)

	// 1. 查缓存
	if val, err := db.RedisClient.Get(cacheKey).Result(); err == nil {
		// 反序列化（这里示意，实际用 json.Unmarshal）
		_ = val
		// return deserialize(val)
	}

	// 2. 查数据库（伪代码）
	// user, err := db.Gain.First(...)
	// if err != nil { return nil, err }

	// 3. 写缓存（设置过期时间，防止永久脏数据）
	// db.RedisClient.Set(cacheKey, serialize(user), 10*time.Minute)

	return nil, nil
}
