package api

import (
	"dbdemo/dao"
	"dbdemo/model"
	"dbdemo/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSimilarQuestions(c *gin.Context) {
	wrongID := c.Query("wrong_id")
	limitStr := c.Query("limit")

	if wrongID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 300,
			"msg":  "wrong_id is required",
		})
		return
	}

	var limit int
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	} else {
		limit = 5 // 默认返回5道题目
	}

	similarQuestions, code := dao.GetSimilarQuestions(wrongID, limit)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": similarQuestions,
		"msg":  "success",
	})
}

func AddRecommendationFeedback(c *gin.Context) {
	var feedback model.RecommendationFeedback
	err := c.ShouldBindJSON(&feedback)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}

	code := dao.AddRecommendationFeedback(feedback)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "feedback received",
	})
}
