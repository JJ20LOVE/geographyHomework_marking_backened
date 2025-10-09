package api

import (
	"dbdemo/dao"
	"dbdemo/model"
	"dbdemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllStudent(c *gin.Context) {
	students, code := dao.GetAllStudent()
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
		"data": students,
	})
}

func DeleteStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := dao.DeleteStudent(id)
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

func AddStudent(c *gin.Context) {
	var s model.BaseStudent
	err := c.ShouldBindJSON(&s)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	code := dao.AddStudent(s)
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

func UpdateStudent(c *gin.Context) {
	var s model.Student
	err := c.ShouldBindJSON(&s)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  utils.GetErrMsg(201),
		})
		return
	}
	code := dao.UpdateStudent(s)
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

func GetStudentByClass(c *gin.Context) {
	classID, _ := strconv.Atoi(c.Query("id"))
	students, code := dao.GetStudentByClass(classID)
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
		"data": students,
	})
}
