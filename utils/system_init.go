package utils

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

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

	fmt.Println("config app init")
}

var DB *gorm.DB

func InitMySQL() {
	// 自定义日志模版，打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	fmt.Println("MySQL inited")

	// 测试代码
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println("\n----------->")
	//fmt.Println(user)
}
