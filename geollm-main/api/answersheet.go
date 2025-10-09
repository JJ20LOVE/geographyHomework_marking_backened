package api

import (
	"dbdemo/dao"
	"dbdemo/model"
	"dbdemo/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
)

var ocrQueue = make(chan ocrRequest, 10) // 定义一个有缓冲的 channel 来限制并发量
var wg sync.WaitGroup

type ocrRequest struct {
	aid         int
	examID      int
	fileHeaders []*multipart.FileHeader
	resultChan  chan int
}

func InitWorker() {
	// 启动固定数量的 worker 来处理 OCR 请求
	for i := 0; i < 2; i++ { // 控制并发 worker 的数量
		go ocrWorker()
	}
}

func ocrWorker() {
	for req := range ocrQueue {
		result := dao.StartOcr(req.aid, req.examID, req.fileHeaders)
		req.resultChan <- result
		wg.Done()
	}
}

func GetAnswerSheetList(c *gin.Context) {
	exam_id := c.Query("exam_id")
	class_id := c.Query("class_id")
	answerSheets, code := dao.GetAnswerSheetList(exam_id, class_id)
	if code != 200 {
		//c.JSON(http.StatusOK, gin.H{
		//	"code": code,
		//	"msg":  utils.GetErrMsg(code),
		//})
		//return
		answerSheets = []model.AnswerSheetWithPic{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200, //code
		"msg":  utils.GetErrMsg(code),
		"data": answerSheets,
	})
}

func DeleteAnswerSheet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := dao.DeleteAnswerSheet(id)
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

func CreateAnswerSheet(c *gin.Context) {
	var Uploader model.AnswerSheetUploader
	err := c.ShouldBind(&Uploader)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}

	files := Uploader.File
	t, code := dao.CheckExamType(Uploader.ExamID)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	if t+1 != len(files) {
		c.JSON(http.StatusOK, gin.H{
			"code": 600,
			"msg":  utils.GetErrMsg(600),
		})
		return
	}

	code, aid := dao.AddAnswerSheet(Uploader.StudentID, Uploader.ExamID)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	resultChan := make(chan int)
	wg.Add(1)
	ocrQueue <- ocrRequest{aid: aid, examID: Uploader.ExamID, fileHeaders: files, resultChan: resultChan}

	// 等待 OCR 处理结果
	code = <-resultChan
	close(resultChan)

	if code != 200 {
		dao.DeleteAnswerSheet(aid)
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}

	err = utils.UploadFiles(files, aid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 203,
			"msg":  utils.GetErrMsg(203),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}

func CorrectOcr(c *gin.Context) {
	var o model.OcrEditor
	err := c.ShouldBindJSON(&o)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	code := dao.CorrectOcr(o)
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

func Evaluator(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := dao.CheckAnswerSheetID(id)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	code = dao.Evaluator(id)
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

func GetAnswerSheetInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := dao.CheckAnswerSheetID(id)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	info, code := dao.GetAnswerSheetInfo(id)
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
		"data": info,
	})
}

func BatchEvaluator(c *gin.Context) {
	exam_id := c.Query("exam_id")
	class_id := c.Query("class_id")
	is_skip, _ := strconv.Atoi(c.Query("is_skip"))
	code := dao.BatchEvaluator(exam_id, class_id, is_skip)
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
	})
}
