package test

import (
	"Gin-Chat/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestGorm(t *testing.T) {

	dsn := "root:Yizhili80@tcp(127.0.0.1:3306)/Gin-Chat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Auto Migrate
	db.AutoMigrate(&models.UserBasic{})
	// Set table options
	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&models.UserBasic{})

	// 插入
	user := &models.UserBasic{Name: "mojo"}
	db.Create(user)

	// 查询
	fmt.Println(db.Find(&models.UserBasic{}, "id = ?", 0))

	// 批量插入
	//var users = []models.UserBasic{user}
	//db.Create(&users)

}