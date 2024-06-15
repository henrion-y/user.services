package dtos

type LoginByPhoneCodeDto struct {
	Phone    string `json:"phone" binding:"required,number,len=11"` // 手机号码
	Code     string `json:"code" binding:"required,number,len=6"`   // 手机验证码
	OauthId  int    `json:"oauth_id"`                               // 第三方登录信息id
	Platform int    `json:"platform" binding:"required"`            // 登录渠道，1:iOS,2:Android,3:Web
}

type LoginByOauthUserCodeDto struct {
	Code string `json:"code" binding:"required"` // 第三方登录唯一标识
}

type LoginByPassWordDto struct {
	Phone    string `json:"phone" binding:"required,number,len=11"` // 手机号码
	Password string `json:"password" binding:"required"`            // 手机验证码
	Platform int    `json:"platform" binding:"required"`            // 登录渠道，1:iOS,2:Android,3:Web
}

type UserInfoDto struct {
	UserId       uint64 `json:"user_id"`       // 用户id， 用来跟其他表做关联
	Phone        string `json:"phone"`         // 手机号
	Username     string `json:"user_name"`     // 昵称
	HeadPortrait string `json:"head_portrait"` // 像头
	UserProfile  string `json:"user_profile"`  // 简介
	Sex          int    `json:"sex"`           // 性别(0保密，1男，2女)
}

type UserLoginDto struct {
	UserInfo UserInfoDto `json:"user_info"`
	JwtToken string      `json:"jwt_token"`
	IsNew    bool        `json:"is_new"`
}

type SetPasswordDto struct {
	UserId      uint64 `json:"user_id"`                         // 用户id， 用来跟其他表做关联
	Password    string `json:"password"`                        // 原密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

type UserIdDto struct {
	UserId uint64 `json:"user_id" form:"user_id" binding:"required"` // 用户id， 用来跟其他表做关联
}

type UserTagDto struct {
	UserId uint64 `json:"user_id" form:"user_id" binding:"required"` // 用户id， 用来跟其他表做关联
	TagId  string `json:"tag_id" form:"tag_id" binding:"required"`
}
