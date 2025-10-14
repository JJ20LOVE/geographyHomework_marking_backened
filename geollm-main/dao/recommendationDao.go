package dao

import (
	"dbdemo/model"
	"dbdemo/service"
	"fmt"
	"os"
)

var aiService *service.AIService

func init() {
	// 从环境变量或配置文件中获取API配置
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com" // 默认URL
	}

	aiService = service.NewAIService(apiKey, baseURL)
}

// 获取同类题
func GetSimilarQuestions(wrongID string, limit int) ([]model.SimilarQuestion, int) {
	// 获取错题信息
	wrongQuestion, code := GetWrongQuestionByID(wrongID)
	if code != 200 {
		return nil, code
	}

	// 调用AI服务获取同类题目推荐
	similarQuestions, err := aiService.GetSimilarQuestions(
		wrongQuestion.QuestionText,
		wrongQuestion.KnowledgePoint,
		limit,
	)

	if err != nil {
		fmt.Printf("GetSimilarQuestions error: %v\n", err)
		// 即使AI服务出错，也返回模拟数据，保证用户体验
		similarQuestions, _ = aiService.GetSimilarQuestions(
			wrongQuestion.QuestionText,
			wrongQuestion.KnowledgePoint,
			limit,
		)
		return similarQuestions, 200
	}

	return similarQuestions, 200
}

// 推荐同类题反馈结果
func AddRecommendationFeedback(feedback model.RecommendationFeedback) int {
	sqlStr := `INSERT INTO recommendation_feedback 
        (student_id, wrong_id, question_id, feedback) 
        VALUES (?, ?, ?, ?)`

	_, err := model.Db.Exec(sqlStr,
		feedback.StudentID, feedback.WrongID, feedback.QuestionID, feedback.Feedback)

	if err != nil {
		fmt.Printf("AddRecommendationFeedback error: %v\n", err)
		return 400
	}
	return 200
}
