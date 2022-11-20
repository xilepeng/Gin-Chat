package service

import (
	"Gin-Chat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 所有用户
// @Tags 获取用户列表
// @Success 200 {string} data
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// CreateUser
// @Tags 用户模块
// @Summary 新增用户
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassWord query string false "确认密码"
// @Success 200 {string} data
// @Router  /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassWord := c.Query("repassWord")
	if password != repassWord {
		c.JSON(-1, gin.H{
			"message": "密码不一致",
		})
	}
	user.PassWord = password
	models.CreateUser(user)
	c.JSON(200, gin.H{"message": "创建用户成功！"})
}
