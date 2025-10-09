package api

import (
	"dbdemo/dao"
	"dbdemo/model"
	"dbdemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllClass(c *gin.Context) {
	classes, code := dao.GetAllClass()
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.GetErrMsg(code),
		"data": classes,
	})
}

func AddClass(c *gin.Context) {
	var cl model.BaseClass
	err := c.ShouldBindJSON(&cl)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	code := dao.AddClass(cl.ClassName)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}

func UpdateClass(c *gin.Context) {
	var cl model.Class
	err := c.ShouldBindJSON(&cl)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  utils.GetErrMsg(201),
		})
		return
	}
	code := dao.UpdateClass(cl.ClassID, cl.ClassName)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}

func DeleteClass(c *gin.Context) {
	classid, _ := strconv.Atoi(c.Query("class_id"))
	code := dao.DeleteClass(classid)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}
