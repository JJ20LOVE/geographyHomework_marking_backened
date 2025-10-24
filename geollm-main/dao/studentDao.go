package dao

import "dbdemo/model"

func GetAllStudent() ([]model.Student, int) {
	var student []model.Student
	err := model.Db.Select(&student, "SELECT * FROM student")
	if err != nil {
		return nil, 400
	}
	if len(student) == 0 {
		return nil, 500
	}
	return student, 200
}

func CheckStudentID(id string) (int, int) {
	sqlStr := "select * from student where student_id = ?"
	var student []model.Student
	err := model.Db.Select(&student, sqlStr, id)
	if err != nil {
		return 0, 400
	}
	if len(student) == 0 {
		return 0, 500
	}
	return student[0].ID, 200
}

func DeleteStudent(id int) int {
	code := CheckID(id)
	if code != 200 {
		return code
	}
	code = CheckAnswerSheetByStudentID(id)
	if code == 200 {
		return 505
	}
	sqlStr := "delete from student where id = ?"
	_, err := model.Db.Exec(sqlStr, id)
	if err != nil {
		return 400
	}
	return 200
}

func AddStudent(student model.BaseStudent) int {
	code := CheckClassID(student.ClassID)
	if code != 200 {
		return 502
	}
	_, code = CheckStudentID(student.StudentID)
	if code == 200 {
		return 501
	}

	sqlStr := "insert into student(student_id, student_name,class_id) values(?, ? , ?)"
	_, err := model.Db.Exec(sqlStr, student.StudentID, student.StudentName, student.ClassID)
	if err != nil {
		return 400
	}
	return 200
}

func UpdateStudent(student model.Student) int {
	code := CheckID(student.ID)
	if code != 200 {
		return code
	}
	code = CheckClassID(student.ClassID)
	if code != 200 {
		return 502
	}
	sqlStr := "update student set student_id=?, student_name = ?, class_id= ? where id = ?"
	_, err := model.Db.Exec(sqlStr, student.StudentID, student.StudentName, student.ClassID, student.ID)
	if err != nil {
		return 400
	}
	return 200
}

func CheckID(id int) int {
	sqlStr := "select * from student where id = ?"
	var student []model.Student
	err := model.Db.Select(&student, sqlStr, id)
	if err != nil {
		return 400
	}
	if len(student) == 0 {
		return 500
	}
	return 200
}

func GetStudentByClass(classID int) ([]model.Student, int) {
	var student []model.Student
	code := CheckClassID(classID)
	if code != 200 {
		return nil, code
	}
	err := model.Db.Select(&student, "SELECT * FROM student WHERE class_id = ?", classID)
	if err != nil {
		return nil, 400
	}
	if len(student) == 0 {
		return nil, 500
	}
	return student, 200
}
func GetStudentById(id int) (model.Student, int) {
	var student model.Student
	sqlStr := "SELECT * FROM student WHERE id = ?"
	err := model.Db.Get(&student, sqlStr, id)
	if err != nil {
		return model.Student{}, 400
	}
	return student, 200
}
