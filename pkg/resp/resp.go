package resp

import (
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
	c.JSON(http.StatusOK, gin.H{
		"code": func() int {
			if len(code) != 0 {
				return code[0]
			}
			return 500
		}(),
		"data": nil,
		"msg": func() string {
			if msg != "" {
				return msg
			}
			return err.Error()
		}(),
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
