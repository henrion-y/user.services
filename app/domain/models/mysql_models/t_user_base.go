package mysql_models

import "time"

type TUserBase struct {
	UserId       uint64     `gorm:"column:user_id;primary_key" json:"user_id"` // 主键 用户id， 用来跟其他表做关联
	UserType     int        `gorm:"column:user_type" json:"user_type"`         // 用户类型，1、普通用户，2、作者，3、认证律师. 100、管理员
	Phone        string     `gorm:"column:phone" json:"phone"`                 // 手机号
	Username     string     `gorm:"column:user_name" json:"user_name"`         // 昵称
	UserProfile  string     `gorm:"column:user_profile" json:"user_profile"`   // 用户简介
	Password     string     `gorm:"column:password" json:"password"`           // 密码
	HeadPortrait string     `gorm:"column:head_portrait" json:"head_portrait"` // 像头
	Sex          int        `gorm:"column:sex" json:"sex"`                     // 性别(0保密，1男，2女)
	UserStatus   int        `gorm:"column:user_status" json:"user_status"`     // 1启用；2 禁用
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`       // 最后更新时间
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`       // 创建时间
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`       // 删除时间
}

const DefaultHeadPortrait = "http://img.alzx.top/img/head_portrait/default.webp"

// TableName sets the insert table name for this struct type
func (u *TUserBase) TableName() string {
	return "t_user_base"
}
