package models

import "gorm.io/gorm"

// 群聊信息
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string // 图片
	Type    int
	Desc    string // 描述

}

func (table *GroupBasic) TableName() string {
	return "group_basic" // 数据库里表的命名
}
