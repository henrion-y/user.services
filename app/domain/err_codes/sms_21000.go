package err_codes

import "github.com/henrion-y/base.services/infra/xerror"

const (
	SmsCodeHasNotExpired = 21000 // 验证码未过期
	SmsCodeSendError     = 21010 // 验证码发送失败
)

var smsErrorMap = xerror.ErrorDefinition{
	SmsCodeHasNotExpired: "验证码未过期",
	SmsCodeSendError:     "验证码发送失败",
}

func init() {
	for k, v := range smsErrorMap {
		xerror.SetAllError(k, v)
	}
}
