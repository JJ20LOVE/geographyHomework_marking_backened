package utils

import (
	"github.com/gin-gonic/gin"
)

func ResponseWithMsg(c *gin.Context, err string) {
	c.JSON(200, gin.H{
		"msg": err,
	})
}
