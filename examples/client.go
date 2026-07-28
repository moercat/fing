//go:build ignore
// +build ignore

// Package main 演示 fing 服务端的 Go 客户端调用方式。
//
// 运行：go run examples/client.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:9765"

type LoginResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type Profile struct {
	ID            uint   `json:"id"`
	UserName      string `json:"user_name"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
	Avatar        string `json:"avatar"`
}

func main() {
	// 1. 登录获取 JWT
	token, err := login("alice", "alice12345")
	if err != nil {
		fmt.Println("登录失败:", err)
		return
	}
	fmt.Println("JWT:", token)

	// 2. 用 JWT 拉取资料
	profile, err := getProfile(token)
	if err != nil {
		fmt.Println("拉取资料失败:", err)
		return
	}
	fmt.Printf("资料: %+v\n", profile)

	// 3. 修改昵称
	if err := updateProfile(token, map[string]string{"nickname": "Alice-New"}); err != nil {
		fmt.Println("修改失败:", err)
		return
	}
	fmt.Println("昵称已更新")
}

func login(user, pass string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"user_name": user,
		"password":  pass,
	})
	resp, err := http.Post(baseURL+"/api/v1/login/jwt", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var lr LoginResp
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return "", err
	}
	if lr.Code != 0 {
		return "", fmt.Errorf("code=%d msg=%s", lr.Code, lr.Msg)
	}
	return lr.Token, nil
}

func getProfile(token string) (*Profile, error) {
	req, _ := http.NewRequest("GET", baseURL+"/api/v2/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("profile resp:", string(body))

	var p Profile
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func updateProfile(token string, fields map[string]string) error {
	body, _ := json.Marshal(fields)
	req, _ := http.NewRequest("PUT", baseURL+"/api/v2/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
