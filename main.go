package main

import (
	"Gin-Chat/router"
	"Gin-Chat/utils"
)

func main() {

	utils.InitConfig()
	utils.InitMySQL()

	r := router.Router()
	err := r.Run()
	if err != nil {
		return
	}
}
