// Package emailtemplate 提供 HTML 邮件模板。
//
// 模板分类：
//   - Welcome         注册欢迎
//   - PasswordReset   密码重置
//   - VerifyEmail     邮箱验证
//   - DailyReport     日报通知
//
// 用法：
//
//	body, _ := emailtemplate.Render(emailtemplate.Welcome, map[string]string{
//	    "Name": "alice",
//	})
//	email.SendMail(to, name, subject, body)
package emailtemplate

import (
	"bytes"
	"fmt"
	"html/template"
)

// TemplateName 邮件模板名
type TemplateName string

const (
	Welcome       TemplateName = "welcome"
	PasswordReset TemplateName = "password_reset"
	VerifyEmail   TemplateName = "verify_email"
	DailyReport   TemplateName = "daily_report"
)

var templates = map[TemplateName]string{
	Welcome: `<h2>欢迎加入 {{.Name}}</h2>
<p>感谢注册，请点击下方链接完成邮箱验证：</p>
<p><a href="{{.VerifyURL}}" style="display:inline-block;padding:10px 20px;background:#4F46E5;color:white;text-decoration:none;border-radius:4px;">验证邮箱</a></p>`,

	PasswordReset: `<h2>密码重置</h2>
<p>您正在重置密码，请点击下方链接完成操作：</p>
<p><a href="{{.ResetURL}}" style="display:inline-block;padding:10px 20px;background:#4F46E5;color:white;text-decoration:none;border-radius:4px;">重置密码</a></p>
<p>链接 30 分钟内有效。如果不是您本人操作，请忽略此邮件。</p>`,

	VerifyEmail: `<h2>邮箱验证</h2>
<p>请点击下方链接验证您的邮箱：</p>
<p><a href="{{.VerifyURL}}">{{.VerifyURL}}</a></p>`,

	DailyReport: `<h2>{{.Title}}</h2>
<p>时间：{{.Date}}</p>
<table border="1" cellpadding="8" style="border-collapse:collapse;">
{{range .Items}}<tr><td>{{.Name}}</td><td>{{.Value}}</td></tr>{{end}}
</table>
<p style="color:gray;font-size:12px;">本邮件由系统自动发送。</p>`,
}

// Render 渲染邮件模板，返回 HTML 字符串
func Render(name TemplateName, data map[string]any) (string, error) {
	tplStr, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("template %q not found", name)
	}
	t, err := template.New(string(name)).Parse(tplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderWithSubject 渲染并返回 (subject, htmlBody)
func RenderWithSubject(name TemplateName, data map[string]any) (subject, body string, err error) {
	subjects := map[TemplateName]string{
		Welcome:       "欢迎加入 fing",
		PasswordReset: "密码重置",
		VerifyEmail:   "邮箱验证",
		DailyReport:   "{{.Title}}",
	}
	s, ok := subjects[name]
	if !ok {
		return "", "", fmt.Errorf("template %q not found", name)
	}
	if t, err := template.New("subject").Parse(s); err == nil {
		var buf bytes.Buffer
		_ = t.Execute(&buf, data)
		subject = buf.String()
	} else {
		subject = s
	}
	body, err = Render(name, data)
	return
}