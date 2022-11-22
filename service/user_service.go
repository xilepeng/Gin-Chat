package service

import (
	"Gin-Chat/models"
	"Gin-Chat/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"

	"github.com/gin-gonic/gin"
)

// CreateUser
// @Tags 用户模块
// @Summary 创建用户
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassWord query string false "确认密码"
// @param phone query string false "phone"
// @param email query string false "email"
// @Success 200 {string} data
// @Router  /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	u := models.FindUserByName(user.Name)
	if u.Name != "" {
		c.JSON(-1, gin.H{
			"message": "该用户名已注册！",
		})
		return
	}
	salt := fmt.Sprintf("%06d", rand.Int31())

	password := c.Query("password")
	repassWord := c.Query("repassWord")
	if password != repassWord {
		c.JSON(-1, gin.H{
			"message": "密码不一致！",
		})
		return
	}
	// user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt) // 加密后存储
	user.Phone = c.Query("phone")
	u = models.FindUserByPhone(user.Phone)
	if u.Phone != "" {
		c.JSON(-1, gin.H{
			"message": "该手机号已注册！",
		})
		return
	}
	user.Email = c.Query("email")
	u = models.FindUserByEmail(user.Email)
	if u.Email != "" {
		c.JSON(-1, gin.H{
			"message": "该邮箱已注册！",
		})
		return
	}

	models.CreateUser(user)
	c.JSON(200, gin.H{"message": "创建用户成功！"})
}

// DeleteUser
// @Tags 用户模块
// @Summary 删除用户
// @param id query string false "id"
// @Success 200 {string} data
// @Router  /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{"message": "删除用户成功！"})
}

// UpdateUser
// @Tags 用户模块
// @Summary 修改用户
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} data
// @Router  /user/updateUser [put]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("update:", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"msg": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{"message": "修改用户成功！"})
	}

}

// GetUserList
// @Tags 用户模块
// @Summary 获取所有用户
// @Success 200 {string} data
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "查找用户成功！",
	})
}

// Login
// @Tags 用户模块
// @Summary 登录
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} data
// @Router /user/Login [get]
func Login(c *gin.Context) {
	user := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	if user.Identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "该用户不存在！",
		})
	}

}
