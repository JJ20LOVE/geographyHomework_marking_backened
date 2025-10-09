package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"smg":  GetErrMsg(200),
		"data": v,
	})
}

// Failed 普通的操作失败返回
func Failed(c *gin.Context, v int) {
	c.JSON(http.StatusOK, gin.H{
		"code": v,
		"msg":  GetErrMsg(v),
	})
}
