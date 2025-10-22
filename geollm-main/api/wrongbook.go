package api

import (
	"dbdemo/dao"
	"dbdemo/model"
	"dbdemo/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddWrongQuestion(c *gin.Context) {
	var wq model.WrongQuestionRequest
	err := c.ShouldBindJSON(&wq)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}

	code := dao.AddWrongQuestion(wq)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "added to wrongbook successfully",
	})
}

func GetWrongQuestionsByStudent(c *gin.Context) {
	studentID := c.Query("student_id")
	knowledgePoint := c.Query("knowledge_point")

	if studentID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 300,
			"msg":  "student_id is required",
		})
		return
	}

	wrongQuestions, code := dao.GetWrongQuestionsByStudent(studentID, knowledgePoint)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": wrongQuestions,
		"msg":  "success",
	})
}

func GetWrongQuestionByID(c *gin.Context) {
	var wq model.WrongQuestion
	wrongID, _ := strconv.Atoi(c.Query("wrong_id"))
	wq, code := dao.GetWrongQuestionByID(wrongID)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": wq,
		"msg":  "success",
	})
}

func DeleteWrongQuestion(c *gin.Context) {
	wrongID := c.Query("wrong_id")

	if wrongID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 300,
			"msg":  "wrong_id is required",
		})
		return
	}

	code := dao.DeleteWrongQuestion(wrongID)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
