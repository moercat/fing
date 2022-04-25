package middleware

import (
	"fing/pkg/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"strings"
)

func Cover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if c.IsAborted() {
				c.Status(200)
			}
			switch errStr := err.(type) {
			case string:
				p := strings.Split(errStr, "#")
				if len(p) == 3 && p[0] == "recover" {
					resp.Fail(c, nil, p[2])
					return
				} else {
					fmt.Println(err)
					debug.PrintStack()
				}
			default:
				fmt.Println(err)
				debug.PrintStack()
			}
		}
	}()
	c.Next()
}
