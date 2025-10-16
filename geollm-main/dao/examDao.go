package dao

import (
	"dbdemo/model"
	"dbdemo/utils"
	"encoding/json"
	"fmt"
	"time"
)

func AddExam(title string, creater, qnumber, t int) (int, int) {
	code := CheckUserID(creater)
	if code != 200 {
		return 0, code
	}
	code = CheckExamTitle(title)
	if code == 200 {
		return 0, 604
	}
	date := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := "INSERT INTO exam(title,create_date, creater,qnumber,type) VALUES(?, ?, ?,?,?)"
	result, err := model.Db.Exec(sqlStr, title, date, creater, qnumber, t)
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, 400
	}
	sqlStr = "INSERT INTO exam_detail(exam_id, qid) VALUES(?, ?)"
	for i := 1; i <= qnumber; i++ {
		_, err = model.Db.Exec(sqlStr, lastInsertID, i)
		if err != nil {
			return 0, 400
		}
	}
	return int(lastInsertID), 200
}

func GetAllExam() ([]model.ExamN, int) {
	var exams []model.ExamN
	err := model.Db.Select(&exams, "SELECT * FROM exam")
	if err != nil {
		fmt.Println(err.Error())
		return nil, 400
	}
	if len(exams) == 0 {
		return nil, 605
	}
	return exams, 200
}

func GetExamDetail(exam_id int) (model.YiTuo, int) {
	code := CheckExamID(exam_id)
	if code != 200 {
		fmt.Printf("CheckExamID failed: %d\n", code)
		return model.YiTuo{}, code
	}

	// 添加调试信息
	fmt.Printf("Getting exam detail for exam_id: %d\n", exam_id)

	var exam model.ExamDetail
	exam.ExamID = exam_id

	// 先检查 exam_detail 表数据
	var detailCount int
	err := model.Db.Get(&detailCount, "SELECT COUNT(*) FROM exam_detail WHERE exam_id=?", exam_id)
	if err != nil {
		fmt.Printf("Error counting exam_detail: %v\n", err)
		return model.YiTuo{}, 400
	}
	fmt.Printf("Found %d records in exam_detail for exam_id: %d\n", detailCount, exam_id)

	err = model.Db.Select(&exam.QDetail, "SELECT qid, point, COALESCE(tihao, '') as tihao FROM exam_detail WHERE exam_id=?", exam_id)
	if err != nil {
		fmt.Printf("Error querying exam_detail: %v\n", err)
		return model.YiTuo{}, 400
	}
	fmt.Printf("Successfully loaded %d question details\n", len(exam.QDetail))

	// 检查其他相关表
	var questionCount, answerCount int
	model.Db.Get(&questionCount, "SELECT COUNT(*) FROM exam_question WHERE exam_id=?", exam_id)
	model.Db.Get(&answerCount, "SELECT COUNT(*) FROM exam_answer WHERE exam_id=?", exam_id)
	fmt.Printf("Question records: %d, Answer records: %d\n", questionCount, answerCount)

	exam.Answer, exam.Question, code = Unmarshaler(exam_id)
	if code != 200 {
		fmt.Printf("Unmarshaler failed with code: %d\n", code)
		return model.YiTuo{}, code
	}

	fmt.Printf("Successfully unmarshaled %d answer sections and %d question sections\n",
		len(exam.Answer), len(exam.Question))
	var data model.YiTuo
	data.ExamID = exam_id
	for i := 0; i < len(exam.Question); i++ {
		var temp struct {
			Title     string `json:"Title"`
			Questions []struct {
				Point   int    `json:"Point" binding:"required"`
				Tihao   string `json:"Number" binding:"required"`
				Content string `json:"Content"`
				Answer  string `json:"Answer"`
			} `json:"Questions"`
		}
		temp.Title = exam.Question[i].Title
		for j := 0; j < len(exam.Question[i].Questions); j++ {
			var temp2 struct {
				Point   int    `json:"Point" binding:"required"`
				Tihao   string `json:"Number" binding:"required"`
				Content string `json:"Content"`
				Answer  string `json:"Answer"`
			}
			temp2.Point = exam.QDetail[j].Point
			temp2.Tihao = exam.QDetail[j].Tihao
			temp2.Content = exam.Question[i].Questions[j].Content
			temp2.Answer = exam.Answer[i].Questions[j].Content
			temp.Questions = append(temp.Questions, temp2)
		}
		data.Data = append(data.Data, temp)
	}
	return data, 200
}

func UpdateExam(id int, title string) int {
	code := CheckExamID(id)
	if code != 200 {
		return code
	}
	sqlStr := "SELECT * FROM exam WHERE title=?"
	var exam []model.ExamN
	err := model.Db.Select(&exam, sqlStr, title)
	if err != nil {
		return 400
	}
	if len(exam) != 0 && exam[0].ExamID != id {
		return 604
	}
	sqlStr = "UPDATE exam SET title = ? WHERE exam_id = ?"
	_, err = model.Db.Exec(sqlStr, title, id)
	if err != nil {
		return 400
	}
	return 200
}

func CheckExamID(id int) int {
	sqlStr := "SELECT * FROM exam WHERE exam_id=?"
	var exam []model.ExamN
	err := model.Db.Select(&exam, sqlStr, id)
	if err != nil {
		return 400
	}
	if len(exam) == 0 {
		return 605
	}
	return 200
}

func CheckExamTitle(title string) int {
	sqlStr := "SELECT * FROM exam WHERE title=?"
	var exam []model.ExamN
	err := model.Db.Select(&exam, sqlStr, title)
	if err != nil {
		return 400
	}
	if len(exam) == 0 {
		return 605
	}
	return 200
}

func DeleteExam(id int) int {
	code := CheckExamID(id)
	if code != 200 {
		return code
	}
	code = CheckAnswerSheetByExamID(id)
	if code == 200 {
		return 607
	}
	_, err := model.Db.Exec("DELETE FROM exam WHERE exam_id=?", id)
	if err != nil {
		return 400
	}
	_, err = model.Db.Exec("DELETE FROM exam_question WHERE exam_id=?", id)
	if err != nil {
		return 400
	}
	_, err = model.Db.Exec("DELETE FROM exam_answer WHERE exam_id=?", id)
	if err != nil {
		return 400
	}
	_, err = model.Db.Exec("DELETE FROM exam_detail WHERE exam_id=?", id)
	if err != nil {
		return 400
	}
	return 200
}

func GetExamTitle(id int) (string, int) {
	sqlStr := "SELECT title FROM exam WHERE exam_id=?"
	var title string
	err := model.Db.Get(&title, sqlStr, id)
	if err != nil {
		return "", 605
	}
	return title, 200
}

func QuestionExtractor(id int) int {
	sections, err := utils.Extractor("./uploads/tmp/exam/question.docx")
	jsonData, err := json.MarshalIndent(sections, "", "    ")
	if err != nil {
		return 205
	}
	sqlStr := "INSERT INTO exam_question(exam_id, question) VALUES(?, ?)"
	_, err = model.Db.Exec(sqlStr, id, jsonData)
	if err != nil {
		return 400
	}
	for _, section := range sections {
		for _, question := range section.Questions {
			sqlStr = "UPDATE exam_detail SET point = 0,tihao=? WHERE exam_id = ? AND qid = ?"
			_, err = model.Db.Exec(sqlStr, question.Number, id, question.Number)
			if err != nil {
				return 400
			}
		}
	}
	return 200
}

func AnswerExtractor(id int) int {
	sections, err := utils.Extractor("./uploads/tmp/exam/answer.docx")
	jsonData, err := json.MarshalIndent(sections, "", "    ")
	if err != nil {
		return 205
	}
	sqlStr := "INSERT INTO exam_answer(exam_id, answer) VALUES(?, ?)"
	_, err = model.Db.Exec(sqlStr, id, jsonData)
	if err != nil {
		return 400
	}
	return 200
}

func DeUploader(id int) int {
	code := CheckExamID(id)
	if code != 200 {
		return code
	}
	sqlStr := "DELETE FROM exam_question WHERE exam_id=?"
	_, err := model.Db.Exec(sqlStr, id)
	if err != nil {
		return 400
	}
	sqlStr = "DELETE FROM exam_answer WHERE exam_id=?"
	_, err = model.Db.Exec(sqlStr, id)
	if err != nil {
		return 400
	}
	return 200
}

func SetQuestionDetail(ep model.QuestionDetail) int {
	code := CheckExamID(ep.ExamID)
	if code != 200 {
		return code
	}
	sqlStr := "UPDATE exam_detail SET point = ?,tihao=? WHERE exam_id = ? AND qid = ?"
	_, err := model.Db.Exec(sqlStr, ep.Point, ep.Tihao, ep.ExamID, ep.QuestionID)
	if err != nil {
		return 400
	}
	return 200
}

//func CorrectQuestionExtractor(id int, jsonData []byte) int {
//	sqlStr := "UPDATE exam_question SET question=? WHERE exam_id = ?"
//	_, err := model.Db.Exec(sqlStr, jsonData, id)
//	if err != nil {
//		return 400
//	}
//	return 200
//}
//
//func CorrectAnswerExtractor(id int, jsonData []byte) int {
//	sqlStr := "UPDATE exam_answer SET answer=? WHERE exam_id = ?"
//	_, err := model.Db.Exec(sqlStr, jsonData, id)
//	if err != nil {
//		fmt.Println(err.Error())
//		return 400
//	}
//	return 200
//}

func CheckExamType(id int) (int, int) {
	sqlStr := "SELECT type FROM exam WHERE exam_id=?"
	var exam model.ExamN
	err := model.Db.Get(&exam, sqlStr, id)
	if err != nil {
		return 0, 400
	}
	return exam.Type, 200
}

func Unmarshaler(exam_id int) ([]utils.Section, []utils.Section, int) {
	sqlStr := "select answer from exam_answer where exam_id = ?"
	var answer string
	err := model.Db.Get(&answer, sqlStr, exam_id)
	if err != nil {
		return nil, nil, 400
	}
	sqlStr = "select question from exam_question where exam_id = ?"
	var question string
	err = model.Db.Get(&question, sqlStr, exam_id)
	if err != nil {
		return nil, nil, 400
	}
	var sections1 []utils.Section
	var sections2 []utils.Section
	err = json.Unmarshal([]byte(answer), &sections1)
	if err != nil {
		return nil, nil, 207
	}
	err = json.Unmarshal([]byte(question), &sections2)
	if err != nil {
		return nil, nil, 207
	}
	return sections1, sections2, 200
}

func Marshaler(exam_id int, sections1, sections2 []utils.Section) int {
	// 将 sections1 转换为 JSON 格式
	answer, err := json.Marshal(sections1)
	if err != nil {
		return 207 // JSON 转换失败
	}

	// 将 sections2 转换为 JSON 格式
	question, err := json.Marshal(sections2)
	if err != nil {
		return 207 // JSON 转换失败
	}

	// 定义 SQL 语句，用于插入数据
	sqlStr1 := "UPDATE exam_answer set answer=? WHERE exam_id = ?"
	sqlStr2 := "UPDATE exam_question set question=? WHERE exam_id = ?"

	// 插入 answer 数据
	_, err = model.Db.Exec(sqlStr1, string(answer), exam_id)
	if err != nil {
		fmt.Println(err.Error())

		return 400 // 插入 answer 数据失败
	}

	// 插入 question 数据
	_, err = model.Db.Exec(sqlStr2, string(question), exam_id)
	if err != nil {
		return 400 // 插入 question 数据失败
	}

	// 返回成功状态码
	return 200
}
