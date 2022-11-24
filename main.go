package main

import (
	"Gin-Chat/router"
	"Gin-Chat/utils"
)

func main() {

	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()

	r := router.Router()
	err := r.Run()
	if err != nil {
		return
	}
}
