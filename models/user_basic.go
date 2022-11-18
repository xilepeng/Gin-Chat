package models

import (
	"Gin-Chat/utils"
	"fmt"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity      string `json:"identity,omitempty"`
	Name          string `json:"name,omitempty"`
	PassWord      string `json:"pass_word,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	ClientIp      string `json:"client_ip,omitempty"`
	ClientPort    string `json:"client_port,omitempty"`
	LoginTime     uint64 `json:"login_time,omitempty"`
	HeartBeatTime uint64 `json:"heart_beat_time,omitempty"`
	LoginOutTime  uint64 `json:"login_out_time,omitempty"`
	IsLoginOut    bool   `json:"is_login_out,omitempty"`
	DeviceInfo    string `json:"device_info,omitempty"`
}

func (table *UserBasic) TableName() string {
	return "user_basic" // 数据库里表的命名
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
