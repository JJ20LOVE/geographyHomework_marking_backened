package dao

import (
	"dbdemo/model"
	"dbdemo/utils"
	"log"
	"math"
)

func GetResultByQuestion(exam_id int) ([]model.QuestionScores, int) {
	code := CheckExamID(exam_id)
	if code != 200 {
		return nil, code
	}
	rows, err := model.Db.Query("SELECT qid, answersheet.student_id,student_name,point FROM answersheet_detail JOIN answersheet ON answersheet_detail.aid = answersheet.id JOIN student ON answersheet.student_id=student.id WHERE exam_id = ? ORDER BY qid, student_id", exam_id)
	if err != nil {
		return nil, 400
	}
	defer rows.Close()

	// 用于存储查询结果的 map
	results := make(map[int][]model.StudentScore)
	var qids []int

	// 遍历查询结果
	for rows.Next() {
		var qid int
		var studentID int
		var studentName string
		var point int

		err := rows.Scan(&qid, &studentID, &studentName, &point)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := results[qid]; !ok {
			qids = append(qids, qid)
		}

		// 将学生的得分信息加入当前 qid 对应的学生列表
		results[qid] = append(results[qid], model.StudentScore{
			StudentID:   studentID,
			StudentName: studentName,
			Score:       point,
		})
	}

	// 检查查询过程中是否有错误
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// 构造最终的 JSON 结构
	var orderedResults []model.QuestionScores
	for _, qid := range qids {
		orderedResults = append(orderedResults, model.QuestionScores{
			QuestionID: qid,
			Scores:     results[qid],
		})
	}
	return orderedResults, 200
}

func GetResultByStudent(exam_id int) ([]model.StudentScores, int) {
	code := CheckExamID(exam_id)
	if code != 200 {
		return nil, code
	}
	// 执行查询
	rows, err := model.Db.Query(`SELECT student_id, qid, point FROM answersheet_detail 
                           JOIN answersheet ON answersheet_detail.aid = answersheet.id 
                           WHERE exam_id = ? 
                           ORDER BY student_id, qid`, exam_id)
	if err != nil {
		return nil, 400
	}
	defer rows.Close()

	// 用于存储查询结果的 map，key 是 student_id
	results := make(map[int][]model.QuestionScore)
	var studentIDs []int
	// 遍历查询结果
	for rows.Next() {
		var studentID int
		var qid int
		var point int

		err := rows.Scan(&studentID, &qid, &point)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := results[studentID]; !ok {
			studentIDs = append(studentIDs, studentID)
		}
		// 将每个学生的每道题目得分加入结果
		results[studentID] = append(results[studentID], model.QuestionScore{
			QuestionID: qid,
			Score:      point,
		})
	}

	// 构造最终的 JSON 结构
	var orderedResults []model.StudentScores
	for _, studentID := range studentIDs {
		orderedResults = append(orderedResults, model.StudentScores{
			StudentID: studentID,
			Scores:    results[studentID],
		})
	}
	return orderedResults, 200
}

func GetQuestionPointRate(exam_id int) ([]float64, int) {
	var QP []float64
	var MP []float64
	code := CheckExamID(exam_id)
	if code != 200 {
		return nil, code
	}
	rows, err := model.Db.Query("SELECT avg(point) FROM answersheet_detail JOIN answersheet ON aid=answersheet.id WHERE exam_id = ? GROUP BY qid", exam_id)
	if err != nil {
		return nil, 400
	}
	defer rows.Close()

	for rows.Next() {
		var avg float64
		err := rows.Scan(&avg)
		if err != nil {
			log.Fatal(err)
		}
		QP = append(QP, avg)
	}

	rows, err = model.Db.Query("SELECT point FROM exam_detail WHERE exam_id = ?", exam_id)
	if err != nil {
		return nil, 400
	}
	defer rows.Close()

	for rows.Next() {
		var point float64
		err := rows.Scan(&point)
		if err != nil {
			log.Fatal(err)
		}
		MP = append(MP, point)
	}
	for i := range QP {
		QP[i] = QP[i] / MP[i]
	}
	return QP, 200
}

func GetNameList(exam_id int) ([]model.StudentList, []model.StudentList, int) {
	code := CheckExamID(exam_id)
	if code != 200 {
		return nil, nil, code
	}
	rows, err := model.Db.Query("SELECT student.id, student.student_id, student_name, total_grade FROM student JOIN answersheet ON student.id = answersheet.student_id WHERE exam_id=? ORDER BY total_grade DESC", exam_id)
	if err != nil {
		return nil, nil, 400
	}
	defer rows.Close()

	var studentList []model.StudentList
	var studentList1 []model.StudentList
	var studentList2 []model.StudentList
	for rows.Next() {
		var student model.StudentList
		err := rows.Scan(&student.ID, &student.StudentID, &student.StudentName, &student.Grade)
		if err != nil {
			log.Fatal(err)
		}
		studentList = append(studentList, student)

	}
	p := int(0.27*float64(len(studentList)) + 0.5)
	studentList1 = studentList[:p]
	studentList2 = studentList[len(studentList)-p:]
	return studentList1, studentList2, 200
}

func GetExamData(exam_id int) (model.ExamData, int) {
	var examData model.ExamData
	code := CheckExamID(exam_id)
	if code != 200 {
		return examData, code
	}
	err := model.Db.Get(&examData, "SELECT MAX(total_grade) AS highest,MIN(total_grade) AS lowest,AVG(total_grade) AS average, COUNT(total_grade) AS count FROM answersheet WHERE exam_id=? GROUP BY exam_id", exam_id)
	if err != nil {
		return model.ExamData{}, 400
	}
	return examData, 200
}

func GetStudentInfo(id int) (model.StudentInfo, int) {
	var studentInfo model.StudentInfo
	err := model.Db.Get(&studentInfo, "SELECT student.id, student.student_id,student_name,class_name, AVG(total_grade) AS avg_grade, MAX(total_grade) AS max_grade, MIN(total_grade) AS min_grade FROM student JOIN answersheet ON student.id = answersheet.student_id NATURAL JOIN class WHERE student.id = ? GROUP BY student.id", id)
	if err != nil {
		return studentInfo, 400
	}
	err = model.Db.Select(&studentInfo.History, "SELECT answersheet.id,title,exam.create_date, total_grade FROM answersheet NATURAL JOIN exam WHERE student_id = ?", id)
	if err != nil {
		return studentInfo, 400
	}
	return studentInfo, 200
}

func SOLO(class_id, exam_id, qid int) (model.SOLO, int) {
	var solo model.SOLO
	code := CheckExamID(exam_id)
	if code != 200 {
		return solo, code
	}
	code = CheckClassID(class_id)
	if code != 200 {
		return solo, code
	}
	sqlStr := "SELECT qnumber FROM exam WHERE exam_id=?"
	var qnumber int
	err := model.Db.Get(&qnumber, sqlStr, exam_id)
	if err != nil {
		return solo, 400
	}
	if qid < 1 || qid > qnumber {
		return solo, 400
	}

	section1, section2, code := Unmarshaler(exam_id)
	if code != 200 {
		return solo, code
	}

	index := 0
	question := ""
	answer := ""
	for i := range section1 {
		for j := range section1[i].Questions {
			if index == qid-1 {
				question = section2[i].Questions[j].Content
				answer = section1[i].Questions[j].Content
				break
			}
			index++
		}
	}

	sqlStr = "SELECT point,comment,structure,student.id,student.student_id,student.student_name FROM answersheet_detail JOIN answersheet ON answersheet_detail.aid=answersheet.id JOIN student ON answersheet.student_id = student.id WHERE exam_id=? AND qid=? AND class_id=? ORDER BY point DESC"
	var tmp []struct {
		Comment     string `json:"comment" db:"comment"`
		Point       int    `json:"point" db:"point"`
		Structure   int    `json:"structure" db:"structure"`
		ID          int    `json:"id" db:"id"`
		StudentName string `json:"student_name" db:"student_name"`
		StudentID   int    `json:"student_id" db:"student_id"`
	}
	err = model.Db.Select(&tmp, sqlStr, exam_id, qid, class_id)
	if err != nil {
		return solo, 400
	}

	totalResponses := len(tmp)
	avg := 0.0
	var comments []string
	var students []model.StudentListByQuestion
	stu_structure := []int{0, 0, 0, 0, 0}
	var stu_list [5][]model.StudentListByQuestion //stu_list :=make([][]model.StudentListByQuestion,5)
	for t := range tmp {
		avg += float64(tmp[t].Point)
		comments = append(comments, tmp[t].Comment)
		students = append(students, model.StudentListByQuestion{
			ID:          tmp[t].ID,
			StudentID:   tmp[t].StudentID,
			StudentName: tmp[t].StudentName,
		})
		stu_structure[max(0, tmp[t].Structure-1)]++
		stu_list[max(0, tmp[t].Structure-1)] = append(stu_list[max(0, tmp[t].Structure-1)], model.StudentListByQuestion{
			ID:          tmp[t].ID,
			StudentID:   tmp[t].StudentID,
			StudentName: tmp[t].StudentName,
		})
	}
	if totalResponses != 0 {
		avg /= float64(totalResponses)
	}

	full := 0
	sqlStr = "SELECT point FROM exam_detail WHERE exam_id=? AND qid=?"
	err = model.Db.Get(&full, sqlStr, exam_id, qid)
	if err != nil {
		return solo, 400
	}

	var level int
	if totalResponses > 0 {
		mean := float64(stu_structure[0]*0+stu_structure[1]*1+stu_structure[2]*2+stu_structure[3]*3+stu_structure[4]*4) / float64(totalResponses)
		variance := 0.0
		for i, count := range stu_structure {
			variance += float64(count) * (float64(i) - mean) * (float64(i) - mean)
		}
		variance /= float64(totalResponses)
		stdDev := math.Sqrt(variance)

		if mean > 2.0+stdDev {
			level = 2 // High
		} else if mean < 2.0-stdDev {
			level = 0 // Low
		} else {
			level = 1 // Medium
		}
	}

	p := int(0.27*float64(totalResponses) + 0.5)
	solo.Average = avg
	solo.Problem = utils.Problem(question, comments, avg, full)
	solo.Knowledge = utils.Knowledge(question, answer)
	solo.CLassID = class_id
	solo.StudentNumber = totalResponses
	solo.Highest = tmp[0].Point
	solo.Lowest = tmp[totalResponses-1].Point
	solo.TopStudent = students[:p]
	solo.BackStudent = students[totalResponses-p:]
	solo.PStudentNumber = stu_structure[0]
	solo.UStudentNumber = stu_structure[1]
	solo.MStudentNumber = stu_structure[2]
	solo.RStudentNumber = stu_structure[3]
	solo.EStudentNumber = stu_structure[4]
	solo.Level = level
	solo.PStudentList = stu_list[0]
	solo.UStudentList = stu_list[1]
	solo.MStudentList = stu_list[2]
	solo.RStudentList = stu_list[3]
	solo.EStudentList = stu_list[4]

	return solo, 200
}
