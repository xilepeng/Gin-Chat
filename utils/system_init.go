package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

var (
	DB  *gorm.DB
	RDB *redis.Client
)

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

func InitRedis() {
	var ctx = context.Background()
	RDB = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"), // no password set
		DB:           viper.GetInt("redis.DB"),          // use default DB
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := RDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis inited err:", err)
	} else {
		fmt.Println("Redis inited ->", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到 Redis
func Publish(ctx context.Context, channel string, message string) error {
	var err error
	fmt.Println("Publish ----->", message)
	err = RDB.Publish(ctx, channel, message).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅消息从 Redis
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RDB.Subscribe(ctx, channel)
	fmt.Println("Subscribe -----> 1", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe -----> 2", msg.Payload)
	return msg.Payload, err
}
