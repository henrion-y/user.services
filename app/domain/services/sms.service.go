package services

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/henrion-y/base.services/infra/redisapi"
	"github.com/henrion-y/base.services/infra/xerror"
	"user.services/app/domain/consts"
	"user.services/app/domain/dtos"
	"user.services/app/domain/err_codes"
	"user.services/pkg/sdks/shansuma/sms"
	"user.services/pkg/utils"
)

type SmsService interface {
	SendCode(c context.Context, dto *dtos.SendCodeDto) error
}

type smsService struct {
	smsClient *sms.Client
	redisApi  *redisapi.RedisApi
}

func NewSmsService(redisApi *redisapi.RedisApi, smsClient *sms.Client) SmsService {
	return &smsService{
		redisApi:  redisApi,
		smsClient: smsClient,
	}
}

func (s smsService) SendCode(c context.Context, dto *dtos.SendCodeDto) error {
	code, err := s.redisApi.Get(c, consts.SEND_SMS_BASE_RDB_KER+dto.Phone, 0)
	if err != nil && err != redis.ErrNil {
		return err
	}
	// todo 这里还要查看剩余过期时间
	if code != "" {
		return xerror.NewXErrorByCode(err_codes.SmsCodeHasNotExpired)
	}
	// 生成验证码
	code, err = utils.GenerateSmsCode()
	if err != nil {
		return err
	}
	// 发送验证码
	err = s.smsClient.SendCode(c, dto.Phone, code)
	if err != nil {
		return xerror.NewXErrorByCode(err_codes.SmsCodeSendError)
	}
	// 设置缓存
	_, err = s.redisApi.SetInterface(c, consts.SEND_SMS_BASE_RDB_KER+dto.Phone, code, 5*60, 0)
	if err != nil {
		return err
	}
	return nil
}
