// Package upload 提供文件上传与下载接口。
//
// 端点：
//   - POST /api/v2/upload            单文件上传（multipart）
//   - POST /api/v2/upload/multi      多文件上传
//   - GET  /api/v2/upload/:filename  下载文件
//
// 文件保存到 cfg.UploadDir，URL 通过 /api/v2/upload/:filename 访问。
//
// TODO:
//   - 文件大小限制（默认 50MB）
//   - 文件类型白名单（mime 校验）
//   - 存储到 S3 / OSS（替换本地目录）
package upload

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RouterUpload struct{}

// Router 注册上传路由
func (r *RouterUpload) Router(router *gin.Engine) {
	g := router.Group("/api/v2/upload")
	g.POST("", r.single)
	g.POST("/multi", r.multi)
	g.GET("/:filename", r.download)
}

const (
	uploadDir   = "./uploads"
	maxFileSize = 50 << 20 // 50MB
)

// single 单文件上传
func (r *RouterUpload) single(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到文件"})
		return
	}
	defer file.Close()

	url, err := saveFile(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"url":      url,
			"filename": filepath.Base(url),
			"size":     header.Size,
		},
		"msg": "ok",
	})
}

// multi 多文件上传
func (r *RouterUpload) multi(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "表单解析失败"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到文件"})
		return
	}

	urls := make([]string, 0, len(files))
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			continue
		}
		url, err := saveFile(file, header.Filename)
		file.Close()
		if err == nil {
			urls = append(urls, url)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"urls": urls},
		"msg":  fmt.Sprintf("上传成功 %d 个文件", len(urls)),
	})
}

// download 下载文件
func (r *RouterUpload) download(c *gin.Context) {
	filename := c.Param("filename")
	// 防路径穿越
	filename = filepath.Base(filename)
	full := filepath.Join(uploadDir, filename)

	if _, err := os.Stat(full); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		return
	}

	// 推测 MIME
	ext := filepath.Ext(filename)
	ctype := mime.TypeByExtension(ext)
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	c.Header("Content-Type", ctype)
	c.File(full)
}

// saveFile 把上传的文件保存到 uploadDir，返回可访问的 URL 路径
func saveFile(src io.Reader, origName string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}

	// 文件名加时间戳防止冲突
	ext := filepath.Ext(origName)
	base := strings.TrimSuffix(filepath.Base(origName), ext)
	safe := sanitize(base)
	filename := fmt.Sprintf("%s_%d%s", safe, time.Now().UnixNano(), ext)

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 限制大小
	if _, err := io.CopyN(dst, src, maxFileSize); err != nil && err != io.EOF {
		return "", err
	}

	return "/api/v2/upload/" + filename, nil
}

// sanitize 清理文件名特殊字符
func sanitize(name string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
}