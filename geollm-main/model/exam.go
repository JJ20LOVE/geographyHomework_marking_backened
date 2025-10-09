package model

import (
	"dbdemo/utils"
	"mime/multipart"
)

type ExamN struct {
	ExamID     int    `json:"exam_id" db:"exam_id"`
	Title      string `json:"title" db:"title"`
	CreateDate string `json:"create_date" db:"create_date"`
	Creater    int    `json:"creater" db:"creater"`
	Qnumber    int    `json:"qnumber" db:"qnumber"`
	Type       int    `json:"type" db:"type"`
}

type ExamCreater struct {
	Title    string                `form:"title" binding:"required"`
	Creater  int                   `form:"creater" binding:"required"`
	Qnumber  int                   `form:"qnumber" binding:"required"`
	Type     *int                  `form:"type" binding:"required"`
	Question *multipart.FileHeader `form:"question" binding:"required"`
	Answer   *multipart.FileHeader `form:"answer" binding:"required"`
}

type ExamUpdate struct {
	ExamID int    `json:"exam_id" db:"exam_id" binding:"required"`
	Title  string `json:"title" db:"title" binding:"required"`
}

type ExamQuestion struct {
	ExamID   int    `json:"exam_id" db:"exam_id"`
	Question []byte `json:"question" db:"question"`
}

type QuestionDetail struct {
	ExamID     int    `json:"exam_id" db:"exam_id"`
	QuestionID int    `json:"qid" db:"qid" binding:"required"`
	Point      int    `json:"point" db:"point" binding:"required"`
	Tihao      string `json:"tihao" db:"tihao" binding:"required"`
}

type ExamDetail struct {
	ExamID   int              `json:"exam_id" db:"exam_id"`
	QDetail  []QuestionDetail `json:"q_detail" db:"q_detail"`
	Answer   []utils.Section  `json:"answer" db:"answer"`
	Question []utils.Section  `json:"question" db:"question"`
}

type YiTuo struct {
	ExamID int `json:"exam_id" binding:"required"`
	Data   []struct {
		Title     string `json:"Title"`
		Questions []struct {
			Point   int    `json:"Point" binding:"required"`
			Tihao   string `json:"Number" binding:"required"`
			Content string `json:"Content"`
			Answer  string `json:"Answer"`
		} `json:"Questions"`
	} `json:"data"`
}
