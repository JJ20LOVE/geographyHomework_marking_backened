package api

import (
	"dbdemo/dao"
	"dbdemo/model"
	"dbdemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func AddExam(c *gin.Context) {
	var exam model.ExamCreater
	err := c.ShouldBind(&exam)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	id, code := dao.AddExam(exam.Title, exam.Creater, exam.Qnumber, *exam.Type)
	if code != 200 {
		dao.DeleteExam(id)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	customSaveDir := "./uploads/tmp/exam"
	filename := filepath.Base("answer.docx")
	fullFilePath := filepath.Join(customSaveDir, filename)
	err = c.SaveUploadedFile(exam.Answer, fullFilePath)
	if err != nil {
		dao.DeleteExam(id)
		c.JSON(http.StatusOK, gin.H{
			"code": 203,
			"msg":  utils.GetErrMsg(203),
		})
		return
	}
	filename = filepath.Base("question.docx")
	fullFilePath = filepath.Join(customSaveDir, filename)
	err = c.SaveUploadedFile(exam.Question, fullFilePath)
	if err != nil {
		dao.DeleteExam(id)
		c.JSON(http.StatusOK, gin.H{
			"code": 203,
			"msg":  utils.GetErrMsg(203),
		})
		return
	}
	code = dao.AnswerExtractor(id)
	if code != 200 {
		dao.DeleteExam(id)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	code = dao.QuestionExtractor(id)
	if code != 200 {
		dao.DeleteExam(id)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	err = os.RemoveAll(customSaveDir)
	if err != nil {
		dao.DeleteExam(id)
		c.JSON(http.StatusOK, gin.H{
			"code": 204,
			"msg":  utils.GetErrMsg(204),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
		"data": id,
	})
}

func GetAllExam(c *gin.Context) {
	exams, code := dao.GetAllExam()
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
		"data": exams,
	})
}

func UpdateExam(c *gin.Context) {
	var exam model.ExamUpdate
	err := c.ShouldBindJSON(&exam)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	code := dao.UpdateExam(exam.ExamID, exam.Title)
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

func DeleteExam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := dao.DeleteExam(id)
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

func DeUploader(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	err := os.RemoveAll("./uploads/tmp/exam")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 204,
			"msg":  utils.GetErrMsg(204),
		})
		return
	}
	code := dao.DeUploader(id)
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

//	func SetQuestionDetail(c *gin.Context) {
//		var ep model.QuestionDetail
//		err := c.ShouldBindJSON(&ep)
//		msg := utils.Validate(err)
//		if msg != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code":      300,
//				"msg":       utils.GetErrMsg(300),
//				"validator": msg,
//			})
//			return
//		}
//		code := dao.SetQuestionDetail(ep)
//		if code != 200 {
//			c.JSON(http.StatusOK, gin.H{
//				"code": code,
//				"msg":  utils.GetErrMsg(code),
//			})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code": 200,
//			"msg":  utils.GetErrMsg(200),
//		})
//	}
//
//	func CorrectQuestionExtractor(c *gin.Context) {
//		id, _ := strconv.Atoi(c.Query("id"))
//		jsondata, err := c.GetRawData()
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 201,
//				"msg":  utils.GetErrMsg(201),
//			})
//			return
//		}
//		code := dao.CorrectQuestionExtractor(id, jsondata)
//		if code != 200 {
//			c.JSON(http.StatusOK, gin.H{
//				"code": code,
//				"msg":  utils.GetErrMsg(code),
//			})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code": 200,
//			"msg":  utils.GetErrMsg(200),
//		})
//	}
//
//	func CorrectAnswerExtractor(c *gin.Context) {
//		id, _ := strconv.Atoi(c.Query("id"))
//		jsondata, err := c.GetRawData()
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 201,
//				"msg":  utils.GetErrMsg(201),
//			})
//			return
//		}
//		code := dao.CorrectAnswerExtractor(id, jsondata)
//		if code != 200 {
//			c.JSON(http.StatusOK, gin.H{
//				"code": code,
//				"msg":  utils.GetErrMsg(code),
//			})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code": 200,
//			"msg":  utils.GetErrMsg(200),
//		})
//	}
func Yituo(c *gin.Context) {
	var data model.YiTuo
	err := c.ShouldBindJSON(&data)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 300,
			"msg":  utils.GetErrMsg(300),
		})
		return
	}

	// 创建题目和答案的 Section 切片
	var questionSections []utils.Section
	var answerSections []utils.Section
	counter := 0
	// 处理每个大题
	for _, section := range data.Data {
		// 构建题目 Section 和答案 Section
		questionSection := utils.Section{Title: section.Title}
		answerSection := utils.Section{Title: section.Title}

		// 遍历每个小题
		for _, q := range section.Questions {
			counter++
			// 构建 QuestionDetail 对象
			questionDetail := model.QuestionDetail{
				ExamID:     data.ExamID,
				QuestionID: counter,
				Point:      q.Point,
				Tihao:      q.Tihao,
			}

			// 设置题目详情
			code := dao.SetQuestionDetail(questionDetail)
			if code != 200 {
				c.JSON(http.StatusOK, gin.H{
					"code": code,
					"msg":  utils.GetErrMsg(code),
				})
				return
			}

			// 构建题目和答案的 Question
			question := utils.Question{
				Number:  q.Tihao,
				Content: q.Content,
			}
			answer := utils.Question{
				Number:  q.Tihao,
				Content: q.Answer,
			}

			// 将小题添加到对应的大题中
			questionSection.Questions = append(questionSection.Questions, question)
			answerSection.Questions = append(answerSection.Questions, answer)
		}

		// 将构建好的大题添加到对应的 Section 切片中
		questionSections = append(questionSections, questionSection)
		answerSections = append(answerSections, answerSection)
	}

	code := dao.Marshaler(data.ExamID, answerSections, questionSections)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
	}
	// 如果所有步骤成功，返回成功状态
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}

func GetExamDetail(c *gin.Context) {
	exam_id, _ := strconv.Atoi(c.Query("id"))
	exam, code := dao.GetExamDetail(exam_id)
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
		"data": exam,
	})
}
