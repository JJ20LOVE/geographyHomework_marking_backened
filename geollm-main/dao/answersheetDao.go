package dao

import (
	"dbdemo/model"
	"dbdemo/utils"
	"fmt"
	"mime/multipart"
)

func CheckAnswerSheetByStudentID(studentId int) int {
	sqlStr := "select * from answersheet where student_id = ?"
	var a []model.AnswerSheet
	err := model.Db.Select(&a, sqlStr, studentId)
	if err != nil {
		return 400
	}
	if len(a) == 0 {
		return 606
	}
	return 200
}

func CheckAnswerSheetByExamID(examId int) int {
	sqlStr := "select * from answersheet where exam_id = ?"
	var a []model.AnswerSheet
	err := model.Db.Select(&a, sqlStr, examId)
	if err != nil {
		return 400
	}
	if len(a) == 0 {
		return 606
	}
	return 200
}

func GetAnswerSheetList(exam_id, class_id string) ([]model.AnswerSheetWithPic, int) {
	sqlStr := "SELECT answersheet.id, student.student_id, exam_id, total_grade,is_eva FROM answersheet JOIN student ON answersheet.student_id = student.id WHERE (exam_id = ? OR ? = '') AND (class_id = ? OR ? ='');"
	var a []model.ASInfo
	err := model.Db.Select(&a, sqlStr, exam_id, exam_id, class_id, class_id)
	if err != nil {
		return nil, 400
	}
	if len(a) == 0 {
		return nil, 606
	}
	var b []model.AnswerSheetWithPic
	for _, v := range a {
		sqlStr = "SELECT type FROM exam WHERE exam_id=?"
		var examType string
		err = model.Db.Get(&examType, sqlStr, v.ExamID)
		if err != nil {
			return nil, 400
		}
		var picUrls []string
		picUrls, _ = utils.GetFileUrl(v.ID, 1, examType)
		b = append(b, model.AnswerSheetWithPic{
			AnswerSheet: v,
			PicUrl:      picUrls[0], //只保留一张图
		})
	}
	return b, 200
}

func DeleteAnswerSheet(id int) int {
	//还要删掉答题纸的相关信息。ocr，小分
	code := CheckAnswerSheetID(id)
	if code != 200 {
		return code
	}
	sqlStr := "select type from exam natural join answersheet where id=?"
	var t int
	err := model.Db.Get(&t, sqlStr, id)
	if err != nil {
		return 400
	}
	sqlStr = "delete from answersheet where id=?"
	_, err = model.Db.Exec(sqlStr, id)
	if err != nil {
		return 400
	}
	sqlStr = "delete from ocr where aid=?"
	_, err = model.Db.Exec(sqlStr, id)
	if err != nil {
		return 400
	}
	sqlStr = "delete from answersheet_detail where aid=?"
	_, err = model.Db.Exec(sqlStr, id)
	if err != nil {
		return 400
	}
	err = utils.DeleteFile(id, t)
	if err != nil {
		return 609
	}
	return 200
}

func CheckAnswerSheetID(id int) int {
	sqlStr := "select * from answersheet where id = ?"
	var a []model.AnswerSheet
	err := model.Db.Select(&a, sqlStr, id)
	if err != nil {
		return 400
	}
	if len(a) == 0 {
		return 606
	}
	return 200
}

func AddAnswerSheet(StudentID string, ExamID int) (int, int) {
	code := CheckExamID(ExamID)
	if code != 200 {
		return code, 0
	}
	id, code := CheckStudentID(StudentID)
	if code != 200 {
		return code, 0
	}
	sqlStr := "select * from answersheet where student_id = ? and exam_id = ?"
	var a []model.AnswerSheet
	err := model.Db.Select(&a, sqlStr, id, ExamID)
	if err != nil {
		fmt.Println(err.Error())
		return 400, 0
	}
	if len(a) != 0 {
		return 601, 0
	}
	sqlStr = "insert into answersheet (student_id, exam_id) values (?, ?)"
	result, err := model.Db.Exec(sqlStr, id, ExamID)
	if err != nil {
		return 400, 0
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		// 错误处理
		return 400, 0
	}
	var num int
	err = model.Db.Get(&num, "select exam.qnumber from exam where exam_id=?", ExamID)
	if err != nil {
		return 400, 0
	}
	sqlStr = "INSERT INTO answersheet_detail (aid, qid) VALUES (?, ?)"
	for i := 1; i <= num; i++ {
		_, err = model.Db.Exec(sqlStr, lastInsertID, i)
		if err != nil {
			return 400, 0
		}
	}
	return 200, int(lastInsertID)
}

func StartOcr(aid, exam_id int, fileHeaders []*multipart.FileHeader) int {
	var data []model.OcrResult
	ocrresult := utils.StartOcr(fileHeaders)
	if ocrresult == nil {
		return 608
	}
	sqlStr := "insert into ocr (aid, qid, result) values (?, ?, ?)"
	var num int
	err := model.Db.Get(&num, "select exam.qnumber from exam where exam_id=?", exam_id)
	if err != nil {
		fmt.Println(err.Error())

		return 400
	}
	for i := 0; i < num; i++ {
		if len(ocrresult) <= i {
			ocrresult = append(ocrresult, "")
		}
		data = append(data, struct {
			ID     int    `json:"id"`
			Result string `json:"result"`
		}{ID: i + 1, Result: ocrresult[i]})
		_, err := model.Db.Exec(sqlStr, aid, i+1, ocrresult[i])
		if err != nil {
			_, err := model.Db.Exec("delete from ocr where aid = ?", aid)
			if err != nil {
				return 400
			}

			return 400
		}
	}
	return 200
}

func CorrectOcr(o model.OcrEditor) int {
	code := CheckAnswerSheetID(o.AID)
	if code != 200 {
		return code
	}
	sqlStr := "update ocr set result = ? where aid = ? and qid=?"
	_, err := model.Db.Exec(sqlStr, o.Result, o.AID, o.QID)
	if err != nil {
		return 400
	}
	return 200
}

func Evaluator(id int) int {
	t := struct {
		Num    int `db:"qnumber"`
		ExamID int `db:"exam_id"`
	}{}
	err := model.Db.Get(&t, "select exam.qnumber,exam.exam_id from exam natural join geollm.answersheet where id=?", id)
	if err != nil {
		fmt.Println(err.Error())
		return 400
	}

	sections1, sections2, code := Unmarshaler(t.ExamID)
	if code != 200 {
		return code
	}
	var points []int
	sqlStr := "select point from answersheet_detail where aid = ?"
	err = model.Db.Select(&points, sqlStr, id)
	var student_answer []string
	sqlStr = "select result from ocr where aid = ?"
	err = model.Db.Select(&student_answer, sqlStr, id)
	if err != nil {
		return 400
	}
	var total int
	index := 0
	for i := range sections1 {
		//fmt.Println("Title:", sections2[i].Title) // 假设 Title 是相同的
		for j := range sections1[i].Questions {
			result := utils.Eva(sections2[i].Title+sections2[i].Questions[j].Content, sections1[i].Questions[j].Content, student_answer[index], points[index])
			sqlStr := "UPDATE answersheet_detail SET point = ?,comment=?,structure=? WHERE aid = ? AND qid = ?"
			//fmt.Printf("%d th Question %s: %s\n", index+1, sections2[i].Questions[j].Number, sections2[i].Questions[j].Content)
			_, err := model.Db.Exec(sqlStr, result.Score, result.Comment, result.Structure, id, index+1)
			total += result.Score
			index++
			if err != nil {
				return 400
			}
		}
	}
	sqlStr = "UPDATE answersheet SET total_grade = ?,is_eva = TRUE WHERE id = ?"
	_, err = model.Db.Exec(sqlStr, total, id)
	if err != nil {
		return 400
	}
	return 200
}

func GetAnswerSheetInfo(aid int) (model.AnswerSheetInfo, int) {
	code := CheckAnswerSheetID(aid)
	if code != 200 {
		return model.AnswerSheetInfo{}, code
	}
	var ASI model.AnswerSheetInfo
	sqlStr := "select answersheet.id,answersheet.student_id,answersheet.exam_id,answersheet.total_grade,answersheet.is_eva,student.student_name from answersheet join geollm.student on answersheet.student_id = student.id where answersheet.id=?"
	err := model.Db.Get(&ASI.AnswerSheet, sqlStr, aid)
	if err != nil {
		return model.AnswerSheetInfo{}, 400
	}
	sqlStr = "SELECT qid, point,result, IFNULL(comment, '') AS comment FROM answersheet_detail NATURAL JOIN ocr WHERE aid = ?"
	err = model.Db.Select(&ASI.Questions, sqlStr, aid)
	if err != nil {
		fmt.Println(err.Error())
		return model.AnswerSheetInfo{}, 400
	}
	sqlStr = "SELECT type FROM exam WHERE exam_id=?"
	var examType string
	err = model.Db.Get(&examType, sqlStr, ASI.AnswerSheet.ExamID)
	if err != nil {
		return model.AnswerSheetInfo{}, 400
	}
	ASI.PicUrls, _ = utils.GetFileUrl(ASI.AnswerSheet.ID, 0, examType)
	return ASI, 200
}

func BatchEvaluator(exam_id, class_id string, is_skip int) int {
	sqlStr := "SELECT id FROM answersheet WHERE (exam_id = ? OR ? = '') AND student_id IN (SELECT id FROM student WHERE (class_id = ? OR ? = ''))"
	if is_skip == 1 {
		sqlStr = "SELECT id FROM answersheet WHERE (exam_id = ? OR ? = '') AND student_id IN (SELECT id FROM student WHERE (class_id = ? OR ? = '')) AND is_eva = FALSE"
	}
	var ids []int
	err := model.Db.Select(&ids, sqlStr, exam_id, exam_id, class_id, class_id)
	if err != nil {
		return 400
	}
	for _, id := range ids {
		code := Evaluator(id)
		if code != 200 {
			return code
		}
	}
	return 200
}
