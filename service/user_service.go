package service

import (
	"Gin-Chat/models"
	"Gin-Chat/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该用户名已注册！",
		})
		return
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	password := c.Query("password")
	repassWord := c.Query("repassWord")
	if password != repassWord {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "密码不一致！",
		})
		return
	}
	// user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt) // 加密后存储
	user.Phone = c.Query("phone")
	u = models.FindUserByPhone(user.Phone)
	if u.Phone != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该手机号已注册！",
		})
		return
	}
	user.Email = c.Query("email")
	u = models.FindUserByEmail(user.Email)
	if u.Email != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该邮箱已注册！",
		})
		return
	}

	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建用户成功！",
	})
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
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除用户成功！",
	})
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
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "修改用户成功！",
		})
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
		"code":    0,
		"message": "查找用户成功！",
		"data":    data,
	})
}

// Login
// @Tags 用户模块
// @Summary 登录
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/login [post]
func Login(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该用户尚未注册！",
			"data":    data,
		})
		return
	}

	password := c.Query("password")
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该用户不存在！",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)

	flag := utils.ValidPassword(password, user.Salt, pwd)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "密码不正确！",
			"data":    data,
		})
		return
	}

	data = models.Login(name, pwd) // pwd 加密后的密码
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功！",
		"data":    data,
	})

}

// 防止跨域伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}
