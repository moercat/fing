// Package paging 提供通用分页参数和响应封装。
//
// 业务代码只需用 Paging{Page, PageSize}，底层自动校验 + offset/limit 转换。
package paging

import (
	"math"

	"gorm.io/gorm"
)

// Paging 分页参数（HTTP 请求体）
type Paging struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Normalize 返回安全的 Page/PageSize（夹紧到 [1, 200]）
func (p Paging) Normalize() (page, size int) {
	if p.Page < 1 {
		page = 1
	} else {
		page = p.Page
	}
	if p.PageSize < 1 {
		size = 20
	} else if p.PageSize > 200 {
		size = 200
	} else {
		size = p.PageSize
	}
	return
}

// Offset 计算 SQL OFFSET
func (p Paging) Offset() int {
	page, _ := p.Normalize()
	return (page - 1) * p.Limit()
}

// Limit 计算 SQL LIMIT
func (p Paging) Limit() int {
	_, size := p.Normalize()
	return size
}

// Result 通用分页响应
type Result[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Pages    int64 `json:"pages"`
}

// Page 自动调 GORM Scope 并返回 Result
// 用法：
//
//	var users []model.UserView
//	res := paging.Page(db.Gain.Model(&usr.UserInfo{}), &users, paging.Paging{Page:1, PageSize:20})
func Page[T any](db *gorm.DB, out *[]T, p Paging) Result[T] {
	page, size := p.Normalize()

	var total int64
	db.Count(&total)

	db.Offset(p.Offset()).Limit(size).Find(out)

	pages := int64(0)
	if total > 0 {
		pages = int64(math.Ceil(float64(total) / float64(size)))
	}

	return Result[T]{
		List:     *out,
		Total:    total,
		Page:     page,
		PageSize: size,
		Pages:    pages,
	}
}