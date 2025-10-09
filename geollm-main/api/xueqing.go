package api

import (
	"dbdemo/dao"
	"dbdemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetStudentInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	studentInfo, code := dao.GetStudentInfo(id)
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
		"data": studentInfo,
	})
}

func GetResultByQuestion(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("exam_id"))
	results, code := dao.GetResultByQuestion(exam_id)
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
		"data": results,
	})
}

func GetResultByStudent(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("exam_id"))
	results, code := dao.GetResultByStudent(exam_id)
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
		"data": results,
	})
}

func GetQuestionPointRate(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("exam_id"))
	results, code := dao.GetQuestionPointRate(exam_id)
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
		"data": results,
	})
}

func GetNameList(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("exam_id"))
	top, back, code := dao.GetNameList(exam_id)
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
		"data": gin.H{
			"top":  top,
			"back": back,
		},
	})
}

func GetExamData(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("exam_id"))
	results, code := dao.GetExamData(exam_id)
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
		"data": results,
	})
}

func SOLO(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("exam_id"))
	qid, _ := strconv.Atoi(c.Query("qid"))
	class_id, _ := strconv.Atoi(c.Query("class_id"))
	result, code := dao.SOLO(class_id, exam_id, qid)
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
		"data": result,
	})
}
