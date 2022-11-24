package models

import "gorm.io/gorm"

type UserMessage struct {
	gorm.Model
	Identity      string
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartBeatTime uint64
	LoginOutTime  uint64
	IsLoginOut    bool
	DeviceInfo    string
	Salt          string
}

func (table *UserMessage) TableName() string {
	return "user_message" // 数据库里表的命名
}
