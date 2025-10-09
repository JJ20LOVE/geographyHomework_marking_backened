package model

type StudentScore struct {
	StudentID   int    `json:"student_id"`
	StudentName string `json:"student_name"`
	Score       int    `json:"score"`
}

type QuestionScores struct {
	QuestionID int            `json:"question_id"`
	Scores     []StudentScore `json:"scores"`
}

type QuestionScore struct {
	QuestionID int `json:"question_id"`
	Score      int `json:"score"`
}

type StudentScores struct {
	StudentID int             `json:"student_id"`
	Scores    []QuestionScore `json:"scores"`
}

type StudentList struct {
	ID          int    `json:"id"`
	StudentID   int    `json:"student_id"`
	StudentName string `json:"student_name"`
	Grade       int    `json:"grade"`
}

type ExamData struct {
	Highest int     `json:"highest" db:"highest"`
	Lowest  int     `json:"lowest" db:"lowest"`
	Average float64 `json:"average" db:"average"`
	Count   int     `json:"count" db:"count"`
}

type SOLO struct {
	CLassID        int                     `json:"class_id"`
	StudentNumber  int                     `json:"student_number"`
	PStudentNumber int                     `json:"p_student_number"`
	UStudentNumber int                     `json:"u_student_number"`
	MStudentNumber int                     `json:"m_student_number"`
	RStudentNumber int                     `json:"r_student_number"`
	EStudentNumber int                     `json:"e_student_number"`
	PStudentList   []StudentListByQuestion `json:"p_student_list"`
	UStudentList   []StudentListByQuestion `json:"u_student_list"`
	MStudentList   []StudentListByQuestion `json:"m_student_list"`
	RStudentList   []StudentListByQuestion `json:"r_student_list"`
	EStudentList   []StudentListByQuestion `json:"e_student_list"`
	Level          int                     `json:"level"`
	Problem        string                  `json:"problem"`
	Knowledge      []string                `json:"knowledge"`
	TopStudent     []StudentListByQuestion `json:"top_student"`
	BackStudent    []StudentListByQuestion `json:"back_student"`
	Highest        int                     `json:"highest"`
	Lowest         int                     `json:"lowest"`
	Average        float64                 `json:"average"`
}

type StudentListByQuestion struct {
	ID          int    `json:"id"`
	StudentID   int    `json:"student_id"`
	StudentName string `json:"student_name"`
}
