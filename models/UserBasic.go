package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity      string
	Name          string
	PassWord      string
	Phone         string
	Email         string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64
	LoginOutTime  uint64
	IsLoginOut    bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic" // 数据库里表的命名
}
