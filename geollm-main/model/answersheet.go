package model

import "mime/multipart"

type AnswerSheet struct {
	ID         int  `db:"id" json:"id" binding:"required"`
	StudentID  int  `db:"student_id" json:"student_id" binding:"required"`
	ExamID     int  `db:"exam_id" json:"exam_id" binding:"required"`
	TotalGrade int  `db:"total_grade" json:"total_grade" binding:"required"`
	IsEva      bool `db:"is_eva" json:"is_eva" binding:"required"`
	//Comment    string  `db:"comment" json:"comment"`
}

type Comment struct {
	ID      int    `json:"id" db:"id" binding:"required"`
	Comment string `json:"comment" db:"comment"`
}

type BaseAnswerSheet struct {
	StudentID int `db:"student_id" json:"student_id" binding:"required"`
	ExamID    int `db:"exam_id" json:"exam_id" binding:"required"`
}

type OcrResult struct {
	ID     int    `json:"id"`
	Result string `json:"result"`
}

type OcrEditor struct {
	AID    int    `json:"aid" db:"aid" binding:"required"`
	QID    int    `json:"qid" db:"qid" binding:"required"`
	Result string `json:"result" db:"result"`
}

type AnswerSheetUploader struct {
	StudentID string                  `form:"student_id" binding:"required"`
	ExamID    int                     `form:"exam_id" binding:"required"`
	File      []*multipart.FileHeader `form:"file" binding:"required"`
}

type AnswerSheetInfo struct {
	AnswerSheet struct {
		AnswerSheet
		StudentName string `json:"student_name" db:"student_name"`
	} `json:"basic_info"`
	Questions []QuestionList `json:"questions"`
	PicUrls   []string       `json:"pic_urls"`
}

type QuestionList struct {
	QuestionID int    `json:"question_id" db:"qid"`
	OcrResult  string `json:"ocr_result" db:"result"`
	Point      int    `json:"point" db:"point"`
	Comment    string `json:"comment" db:"comment"`
}

type AnswerSheetWithPic struct {
	AnswerSheet ASInfo `json:"basic_info"`
	PicUrl      string `json:"pic_url"`
}

type ASInfo struct {
	ID         int    `db:"id" json:"id" binding:"required"`
	StudentID  string `db:"student_id" json:"student_id" binding:"required"`
	ExamID     int    `db:"exam_id" json:"exam_id" binding:"required"`
	TotalGrade int    `db:"total_grade" json:"total_grade" binding:"required"`
	IsEva      bool   `db:"is_eva" json:"is_eva" binding:"required"`
}
