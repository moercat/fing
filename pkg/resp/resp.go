package resp

import (
	"fing/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type pageResp struct {
	List  interface{} `json:"list"`
	Count int         `json:"count"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// 失败数据处理
func Fail(c *gin.Context, err error, msg string, code ...int) {
	responseCode := 500
	if len(code) != 0 {
		responseCode = code[0]
	}

	// 如果是APIError类型，使用其中的代码
	if apiErr, ok := err.(*errors.APIError); ok {
		responseCode = apiErr.Code
	}

	c.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": nil,
		"msg": func() string {
			if msg != "" {
				return msg
			}
			if err != nil {
				return err.Error()
			}
			return "未知错误"
		}(),
	})
}

// 失败数据处理（基于APIError）
func FailWithError(c *gin.Context, apiErr *errors.APIError) {
	c.JSON(http.StatusOK, gin.H{
		"code": apiErr.Code,
		"data": nil,
		"msg":  apiErr.Message,
	})
}

// 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg ...string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg": func() string {
			if len(msg) != 0 {
				return msg[0]
			}
			return ""
		}(),
	})
}

// 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int, page int, size int, msg ...string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": pageResp{
			result,
			count,
			page,
			size,
		},
		"msg": func() string {
			if len(msg) != 0 {
				return msg[0]
			}
			return ""
		}(),
	})
}

func CkFail(code int, err error, msg ...string) {
	if err != nil {
		cd := strconv.Itoa(code)
		if len(msg) > 0 && msg[0] != "" {
			panic("recover#" + cd + "#" + msg[0])
		} else {
			panic("recover#" + cd + "#" + err.Error())
		}
	}
}
