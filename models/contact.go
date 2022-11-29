package models

import "gorm.io/gorm"

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId uint //谁的关系信息
	Type    int  // 对应的类型
	Desc    string
}

func (table *Contact) TableName() string {
	return "contact" // 数据库里表的命名
}
