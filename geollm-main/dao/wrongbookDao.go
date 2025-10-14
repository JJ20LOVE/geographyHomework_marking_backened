package dao

import (
	"dbdemo/model"
	"fmt"
)

// 添加错题记录
func AddWrongQuestion(wq model.WrongQuestionRequest) int {

	//检查是否已存在相同的错题记录
	sqlStr := "SELECT wrong_id FROM wrongbook WHERE student_id = ? AND exam_id = ? AND question_id = ?"
	var existingID int
	err := model.Db.Get(&existingID, sqlStr, wq.StudentID, wq.ExamID, wq.QuestionID)
	if err == nil {
		//记录已存在，更新它
		return UpdateWrongQuestion(existingID, wq)
	}

	//插入新记录
	sqlStr = "INSERT INTO wrongbook (student_id, exam_id, question_id, question_text, student_answer, correct_answer, analysis, knowledge_point) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = model.Db.Exec(sqlStr, wq.StudentID, wq.ExamID, wq.QuestionID, wq.QuestionText, wq.StudentAnswer, wq.CorrectAnswer, wq.Analysis, wq.KnowledgePoint)
	if err != nil {
		fmt.Println("AddWrongQuestion error:%v\n", err)
		return 400
	}
	return 200
}

// 更新错题记录
func UpdateWrongQuestion(id int, wq model.WrongQuestionRequest) int {
	sqlStr := "UPDATE wrongbook SET question_text = ?, student_answer = ?, correct_answer = ?, analysis = ?, knowledge_point = ? WHERE wrong_id = ?"
	_, err := model.Db.Exec(sqlStr, wq.QuestionText, wq.StudentAnswer, wq.CorrectAnswer, wq.Analysis, wq.KnowledgePoint, id)
	if err != nil {
		fmt.Println("UpdateWrongQuestion error:%v\n", err)
		return 400
	}
	return 200
}

func GetWrongQuestionsByStudent(studentID string, knowledgePoint string) ([]model.WrongQuestion, int) {
	var wrongQuestions []model.WrongQuestion
	var sqlStr string
	var err error
	if knowledgePoint == "" {
		sqlStr = `SELECT wrong_id, question_text, student_answer, correct_answer, 
                         knowledge_point, create_time 
                  FROM wrongbook 
                  WHERE student_id = ? AND knowledge_point = ? 
                  ORDER BY create_time DESC`
		err = model.Db.Select(&wrongQuestions, sqlStr, studentID, knowledgePoint)
	} else {
		sqlStr = `SELECT wrong_id, question_text, student_answer, correct_answer, 
                         knowledge_point, create_time 
                  FROM wrongbook 
                  WHERE student_id = ? 
                  ORDER BY create_time DESC`
		err = model.Db.Select(&wrongQuestions, sqlStr, studentID)
	}

	if err != nil {
		fmt.Println("GetWrongQuestionByStudent error:%v\n", err)
		return wrongQuestions, 400
	}
	if len(wrongQuestions) == 0 {
		return []model.WrongQuestion{}, 200
	}
	return wrongQuestions, 200
}

// 删除错题记录
func DeleteWrongQuestion(wrongID string) int {
	sqlStr := "DELETE FROM wrongbook WHERE wrong_id = ?"
	result, err := model.Db.Exec(sqlStr, wrongID)

	if err != nil {
		fmt.Printf("DeleteWrongQuestion error: %v\n", err)
		return 400
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return 606 // 记录不存在
	}

	return 200
}

// 通过ID获取错题记录
func GetWrongQuestionByID(wrongID string) (model.WrongQuestion, int) {
	var wrongQuestion model.WrongQuestion
	sqlStr := "SELECT * FROM wrongbook WHERE wrong_id = ?"

	err := model.Db.Get(&wrongQuestion, sqlStr, wrongID)
	if err != nil {
		return model.WrongQuestion{}, 606 // 记录不存在
	}

	return wrongQuestion, 200
}
