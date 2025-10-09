package dao

import (
	"dbdemo/model"
)

func GetAllClass() ([]model.Class, int) {
	var classes []model.Class
	err := model.Db.Select(&classes, "SELECT * FROM class")

	if err != nil {
		return nil, 400
	}
	if len(classes) == 0 {
		return nil, 502
	}
	return classes, 200
}

func AddClass(classname string) int {
	code := CheckClass(classname)
	if code == 200 {
		return 503
	}
	_, err := model.Db.Exec("INSERT INTO class (class_name) VALUES (?)", classname)
	if err != nil {
		return 400
	}
	return 200
}

func CheckClass(classname string) int {
	sqlStr := "SELECT * FROM class WHERE class_name=?"
	var class model.Class
	err := model.Db.Get(&class, sqlStr, classname)
	if err != nil {
		return 400
	}
	return 200
}

func UpdateClass(classid int, classname string) int {
	code := CheckClassID(classid)
	if code != 200 {
		return 502
	}
	code = CheckClass(classname)
	if code == 200 {
		return 503
	}
	_, err := model.Db.Exec("UPDATE class SET class_name=? WHERE class_id=?", classname, classid)
	if err != nil {
		return 400
	}
	return 200
}

func CheckClassID(classid int) int {
	sqlStr := "SELECT * FROM class WHERE class_id=?"
	var class []model.Class
	err := model.Db.Select(&class, sqlStr, classid)
	if err != nil {
		return 400
	}
	if len(class) == 0 {
		return 502
	}
	return 200
}

func CheckStudent(classid int) int {
	sqlStr := "SELECT * FROM student WHERE class_id=?"
	var s model.Student
	err := model.Db.Get(&s, sqlStr, classid)
	if err != nil {
		return 400
	}
	return 200
}

func DeleteClass(classid int) int {
	code := CheckClassID(classid)
	if code != 200 {
		return 502
	}
	code = CheckStudent(classid)
	if code == 200 {
		return 504
	}
	_, err := model.Db.Exec("DELETE FROM class WHERE class_id=?", classid)
	if err != nil {
		return 400
	}
	return 200
}
