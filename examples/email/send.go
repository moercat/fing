//go:build ignore
// +build ignore

// Package main 演示 fing 的邮件发送用法。
//
// 用法：在 service 层调用 fing/pkg/email.SendMail，传入收件人、主题、HTML 内容。
package main

import (
	"fmt"

	"fing/pkg/email"
)

// ============================================================
// 基础：发送纯文本 / HTML 邮件
// ============================================================

// exampleText 发送纯文本邮件
func exampleText() error {
	body := "这是一封纯文本邮件的内容"
	return email.SendMail(
		"alice@example.com",  // 收件人邮箱
		"Alice",              // 收件人姓名
		"测试邮件",             // 主题
		body,                 // 内容（自动按 text/html 发送）
	)
}

// ============================================================
// 进阶：发送 HTML 邮件
// ============================================================

// exampleHTML 发送 HTML 邮件（密码重置邮件典型场景）
func exampleHTML() error {
	token := "abc123def456"
	resetLink := fmt.Sprintf("https://your-domain.com/reset?token=%s", token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body>
  <h2>密码重置</h2>
  <p>您正在重置密码，请点击下面的链接完成操作：</p>
  <p><a href="%s" style="display:inline-block;padding:10px 20px;background:#4F46E5;color:white;text-decoration:none;border-radius:4px;">重置密码</a></p>
  <p>链接 30 分钟内有效。如果不是您本人操作，请忽略此邮件。</p>
  <hr>
  <p style="color:gray;font-size:12px;">本邮件由系统自动发送，请勿回复。</p>
</body>
</html>
`, resetLink)

	return email.SendMail(
		"alice@example.com",
		"Alice",
		"密码重置 - fing",
		body,
	)
}

// ============================================================
// 进阶：批量发送（带错误处理）
// ============================================================

// exampleBatch 批量给用户发通知
func exampleBatch(emails []string) {
	for _, to := range emails {
		if err := email.SendMail(to, "User", "系统通知", "<p>...</p>"); err != nil {
			// 单个失败不中断，记录日志
			fmt.Printf("send to %s failed: %v\n", to, err)
		}
	}
}

// ============================================================
// 进阶：模板渲染（推荐做法）
// ============================================================

// exampleTemplate 实际项目建议用 html/template 渲染邮件内容
// 这里只示意，不引入额外依赖
//
// import "html/template"
//
// const emailTpl = `<h1>Hello {{.Name}}</h1>`
//
// func renderEmail(name string) (string, error) {
//     t := template.Must(template.New("mail").Parse(emailTpl))
//     var buf bytes.Buffer
//     if err := t.Execute(&buf, map[string]string{"Name": name}); err != nil {
//         return "", err
//     }
//     return buf.String(), nil
// }
