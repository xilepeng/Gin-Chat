package router

import (
	"Gin-Chat/docs"
	"Gin-Chat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "" // alt + shift + enter 选择包导入

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/index", service.GetIndex)

	r.POST("/user/createUser", service.CreateUser)
	r.GET("/user/getUserList", service.GetUserList)

	return r
}
