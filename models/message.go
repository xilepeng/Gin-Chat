package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromId   string // 发送者
	TargetId string // 接收者
	Type     string // 消息类型：群聊、私聊
	Media    int    //消息类型：文字、图片、音频
	Content  string // 消息体
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message" // 数据库里表的命名
}
