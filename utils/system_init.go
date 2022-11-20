package utils

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// fmt.Println("\n----------->")
	// fmt.Println("MySQL config: ", viper.Get("mysql"))
	//fmt.Println("app config: ", viper.Get("app"))

	fmt.Println("config app inited")
}

var DB *gorm.DB

func InitMySQL() {
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	fmt.Println("MySQL inited")

	// 测试代码
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println("\n----------->")
	//fmt.Println(user)
}
