package controllers

import (
	"github.com/gin-gonic/gin"
	"user.services/app/domain/dtos"
	"user.services/app/domain/services"
)

type SmsController struct {
	smsService services.SmsService
}

func NewSmsController(smsService services.SmsService) *SmsController {
	return &SmsController{
		smsService: smsService,
	}
}

// SendCode
// @Tags 短信
// @Summary 发送短信验证码
// @Description 发送短信验证码
// @Id SendCode
// @Produce  json
// @Param        SendCodeDto  body  dtos.SendCodeDto  true  "发短信"
// @Success 200 {object} ResponseData
// @Router /sms/send_code [post]
func (h *SmsController) SendCode(c *gin.Context) {
	req := &dtos.SendCodeDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}
	err := h.smsService.SendCode(c, req)
	if err != nil {
		responseError(c, err)
		return
	}

	responseSuccess(c)
}
