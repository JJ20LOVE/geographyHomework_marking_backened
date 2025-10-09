package api

import (
	"dbdemo/dao"
	"dbdemo/middleware"
	"dbdemo/model"
	"dbdemo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(c *gin.Context) {
	var u model.SignUpParam
	err := c.ShouldBindJSON(&u)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	code := dao.SignUp(u.Name, u.Password, u.Email)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}

func Login(c *gin.Context) {
	var data model.LoginParam
	err := c.ShouldBindJSON(&data)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 300,
			"msg":  msg,
		})
		return
	}
	code := dao.CheckLogin(data.Username, data.Password)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	token, msg := middleware.GenerateJWT(data.Username)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 322,
			"msg":  utils.GetErrMsg(322),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   utils.GetErrMsg(200),
		"token": token,
	})
}

func ChangePassword(c *gin.Context) {
	var data model.ChangePassParam
	err := c.ShouldBindJSON(&data)
	msg := utils.Validate(err)
	if msg != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":      300,
			"msg":       utils.GetErrMsg(300),
			"validator": msg,
		})
		return
	}
	code := dao.ChangePassword(data.Username, data.OldPass, data.NewPass)
	if code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.GetErrMsg(200),
	})
}
