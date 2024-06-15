package mysql_models

import "time"

type TOauthUser struct {
	Id           uint64    `gorm:"column:id;primary_key" json:"id"`           // 主键
	UserId       uint64    `gorm:"column:user_id" json:"user_id"`             // 用户id， 用来跟其他表做关联
	OauthUserId  string    `gorm:"column:oauth_user_id" json:"oauth_user_id"` // 第三方临时账号用户id
	Code         string    `gorm:"column:code" json:"code"`                   // 第三方登录唯一标识符
	Username     string    `gorm:"column:user_name" json:"user_name"`         // 昵称
	HeadPortrait string    `gorm:"column:head_portrait" json:"head_portrait"` // 像头
	UserStatus   int       `gorm:"column:user_status" json:"user_status"`     // 1启用；2 禁用
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`       // 最后更新时间
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`       // 创建时间
}

// TableName sets the insert table name for this struct type
func (u *TOauthUser) TableName() string {
	return "t_oauth_user"
}
