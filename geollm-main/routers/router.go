package routers

import (
	"dbdemo/api"
	"dbdemo/middleware"
	"dbdemo/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)

	//r := gin.New()
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	ug := r.Group("api/user")
	{
		// User module routing
		ug.POST("/signup", api.SignUp)
		ug.POST("/login", api.Login)
		ug.PUT("/changePass", api.ChangePassword)
	}

	auth := r.Group("/api")
	auth.Use(middleware.JwtToken())

	sg := auth.Group("/student")
	{
		sg.GET("/getAllStudent", api.GetAllStudent)
		sg.GET("/getStudentInfo", api.GetStudentInfo)
		sg.DELETE("/deleteStudent", api.DeleteStudent)
		sg.POST("/addStudent", api.AddStudent)
		sg.PUT("/updateStudent", api.UpdateStudent)
		sg.GET("/getStudentByClass", api.GetStudentByClass)
		sg.GET("/getStudentById", api.GetStudentById)
	}

	ag := auth.Group("/answersheet")
	{
		ag.GET("/getAnswerSheet", api.GetAnswerSheetList)
		ag.DELETE("/deleteAnswerSheet", api.DeleteAnswerSheet)
		ag.POST("/createAnswerSheet", api.CreateAnswerSheet)
		ag.PUT("/correctOcr", api.CorrectOcr)
		ag.GET("/evaluator", api.Evaluator)
		ag.GET("/getAnswerSheetInfo", api.GetAnswerSheetInfo)
		ag.GET("/batchEvaluator", api.BatchEvaluator)
	}

	eg := auth.Group("/exam")
	{
		eg.POST("/addExam", api.AddExam)
		eg.GET("/getAllExam", api.GetAllExam)
		eg.PUT("/updateExam", api.UpdateExam)
		eg.DELETE("/deleteExam", api.DeleteExam)
		eg.DELETE("/deUploader", api.DeUploader)
		//eg.POST("/setQuestionDetail", api.SetQuestionDetail)
		//eg.PUT("/correctQuestionExtractor", api.CorrectQuestionExtractor)
		//eg.PUT("/correctAnswerExtractor", api.CorrectAnswerExtractor)
		eg.PUT("/yituo", api.Yituo)
		eg.GET("/getExamDetail", api.GetExamDetail)
	}

	cg := auth.Group("/class")
	{
		cg.GET("/getAllClass", api.GetAllClass)
		cg.POST("/addClass", api.AddClass)
		cg.PUT("/updateClass", api.UpdateClass)
		cg.DELETE("/deleteClass", api.DeleteClass)
	}

	xg := auth.Group("/xueqing")
	{
		xg.GET("/getResultByQuestion", api.GetResultByQuestion)
		xg.GET("/getResultByStudent", api.GetResultByStudent)
		xg.GET("/getQuestionPointRate", api.GetQuestionPointRate)
		xg.GET("/getNameList", api.GetNameList)
		xg.GET("/getExamData", api.GetExamData)
		xg.GET("/solo", api.SOLO)
	}

	//错题本相关路由
	wg := auth.Group("/wrongbook")
	{
		wg.POST("/addWrongQuestion", api.AddWrongQuestion)
		wg.GET("/getByStudent", api.GetWrongQuestionsByStudent)
		wg.DELETE("/deleteWrongQuestion", api.DeleteWrongQuestion)
		wg.GET("/getById", api.GetWrongQuestionByID)
	}

	//推荐系统相关路由
	rg := auth.Group("/recommendation")
	{
		rg.GET("/getSimilarQuestions", api.GetSimilarQuestions)
		rg.POST("/feedback", api.AddRecommendationFeedback)
	}

	var HttpPort string
	if utils.AppMode == "debug" {
		HttpPort = "127.0.0.1" + utils.HttpPort
	} else {
		HttpPort = utils.HttpPort
	}
	err := r.Run(HttpPort)
	if err != nil {
		return
	}
}
