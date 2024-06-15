package dtos

type SendCodeDto struct {
	Phone string `json:"phone" binding:"required,number,len=11"` // 手机号码
}
