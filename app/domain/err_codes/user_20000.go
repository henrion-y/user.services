package err_codes

import "github.com/henrion-y/base.services/infra/xerror"

// 业务错误码模块
const (
	UserNotRegistered             = 20000 // 用户未注册
	UserPasswordIncorrect         = 20010 // 原密码错误
	UserVerificationCodeIncorrect = 20020 // 验证码错误
	UserPasswordOrPhoneIncorrect  = 10030 // 手机号或密码错误
)

var userErrorMap = xerror.ErrorDefinition{
	UserNotRegistered:             "用户未注册",
	UserPasswordIncorrect:         "原密码错误",
	UserVerificationCodeIncorrect: "验证码错误",
	UserPasswordOrPhoneIncorrect:  "手机号或密码错误",
}

func init() {
	for k, v := range userErrorMap {
		xerror.SetAllError(k, v)
	}
}
