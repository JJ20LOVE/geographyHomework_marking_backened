package dao

import (
	"dbdemo/model"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"log"
)

func CheckLogin(username string, password string) int {
	sqlStr := "select * from user where username = ?"
	var user []model.User
	err := model.Db.Select(&user, sqlStr, username)
	if err != nil {
		return 400
	}
	if len(user) == 0 {
		return 310
	}
	if user[0].Password != ScryptPw(password) {
		return 311
	}
	return 200
}

func SignUp(username, password, email string) int {
	code := CheckUser(username)
	if code != 200 {
		return code
	}
	sqlStr := "insert into user (username,password,email) values(?,?,?)"
	_, err := model.Db.Exec(sqlStr, username, ScryptPw(password), email)
	if err != nil {
		return 400
	}
	return 200
}

func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 43, 54, 65, 76, 87, 98}
	key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	str := base64.StdEncoding.EncodeToString(key)
	return str
}

func CheckUser(username string) int {
	var user []model.User
	sqlStr := "select * from user where username = ?"
	err := model.Db.Select(&user, sqlStr, username)
	if err != nil {
		return 400
	}
	if len(user) > 0 {
		return 301
	}
	return 200
}

func ChangePassword(username, oldPass, newPass string) int {
	code := CheckLogin(username, oldPass)
	if code != 200 {
		return code
	}
	sqlStr := "update user set password = ? where username = ?"
	_, err := model.Db.Exec(sqlStr, ScryptPw(newPass), username)
	if err != nil {
		return 400
	}
	return 200
}

func CheckUserID(UserID int) int {
	sqlStr := "select * from user where user_id = ?"
	var user []model.User
	err := model.Db.Select(&user, sqlStr, UserID)
	if err != nil {
		return 400
	}
	if len(user) == 0 {
		return 310
	}
	return 200
}
