package model

import "time"

type WrongQuestion struct {
	WrongID        int       `db:"wrong_id" json:"wrong_id"`
	StudentID      int       `db:"student_id" json:"student_id"`
	ExamID         int       `db:"exam_id" json:"exam_id"`
	QuestionID     int       `db:"question_id" json:"question_id"`
	QuestionText   string    `db:"question_text" json:"question_text"`
	StudentAnswer  string    `db:"student_answer" json:"student_answer"`
	CorrectAnswer  string    `db:"correct_answer" json:"correct_answer"`
	Analysis       string    `db:"analysis" json:"analysis"`
	KnowledgePoint string    `db:"knowledge_point" json:"knowledge_point"`
	CreateTime     time.Time `db:"create_time" json:"create_time"`
}

type WrongQuestionRequest struct {
	StudentID      int    `json:"student_id" binding:"required"`
	ExamID         int    `json:"exam_id" binding:"required"`
	QuestionID     int    `json:"question_id" binding:"required"`
	QuestionText   string `json:"question_text" binding:"required"`
	StudentAnswer  string `json:"student_answer" binding:"required"`
	CorrectAnswer  string `json:"correct_answer" binding:"required"`
	Analysis       string `json:"analysis"`
	KnowledgePoint string `json:"knowledge_point"`
}

type RecommendationFeedback struct {
	StudentID  int    `json:"student_id" binding:"required"`
	WrongID    int    `json:"wrong_id" binding:"required"`
	QuestionID int    `json:"question_id" binding:"required"`
	Feedback   string `json:"feedback" binding:"required"`
}

type SimilarQuestion struct {
	QuestionID     int    `json:"question_id"`
	QuestionText   string `json:"question_text"`
	KnowledgePoint string `json:"knowledge_point"`
	Difficulty     string `json:"difficulty"`
}
