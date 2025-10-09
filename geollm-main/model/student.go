package model

type Student struct {
	ID          int    `db:"id" json:"id" binding:"required"`
	StudentID   string `db:"student_id" json:"student_id" binding:"required"`
	StudentName string `db:"student_name" json:"student_name" binding:"required"`
	ClassID     int    `db:"class_id" json:"class_id" binding:"required"`
}

type BaseStudent struct {
	StudentID   string `db:"student_id" json:"student_id" binding:"required"`
	StudentName string `db:"student_name" json:"student_name" binding:"required"`
	ClassID     int    `db:"class_id" json:"class_id" binding:"required"`
}

type History struct {
	ID         int     `json:"id" db:"id"`
	Title      string  `json:"title" db:"title"`
	Date       string  `json:"date" db:"create_date"`
	TotalGrade float64 `json:"total_grade" db:"total_grade"`
}

type StudentInfo struct {
	ID          int     `db:"id" json:"id" binding:"required"`
	StudentID   string  `db:"student_id" json:"student_id" binding:"required"`
	StudentName string  `db:"student_name" json:"student_name" binding:"required"`
	ClassName   string  `db:"class_name" json:"class_name"`
	AvgGrade    float64 `db:"avg_grade" json:"avg_grade"`
	MaxGrade    float64 `db:"max_grade" json:"max_grade"`
	MinGrade    float64 `db:"min_grade" json:"min_grade"`
	History     []History
}
