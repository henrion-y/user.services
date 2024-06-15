package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/henrion-y/base.services/infra/redisapi"
	"github.com/henrion-y/base.services/infra/xerror"
	"github.com/henrion-y/base.services/infra/zlog"
	"go.uber.org/zap"
	"user.services/app/domain/err_codes"

	"github.com/henrion-y/base.services/infra/jwt"
	"golang.org/x/crypto/bcrypt"
	"user.services/app/domain/consts"
	"user.services/app/domain/dtos"
	"user.services/app/domain/models/mysql_models"
	"user.services/app/domain/repositories"
)

type UserService interface {
	SetPassword(c context.Context, dto *dtos.SetPasswordDto) error
	EditUserInfo(c context.Context, dto *dtos.UserInfoDto) error
	GetUserInfoByUserId(c context.Context, dto *dtos.UserIdDto) (*dtos.UserInfoDto, error)
	LoginByPassword(c context.Context, dto *dtos.LoginByPassWordDto) (*dtos.UserLoginDto, error)
	LoginByPhoneCode(c context.Context, dto *dtos.LoginByPhoneCodeDto) (*dtos.UserLoginDto, error)
	LoginByOauthUserCode(c context.Context, dto *dtos.LoginByOauthUserCodeDto) (*dtos.UserLoginDto, error)
}

type userService struct {
	redisApi       *redisapi.RedisApi
	authService    jwt.AuthService
	userRepository repositories.UserRepository
}

func NewUserService(redisApi *redisapi.RedisApi,
	authService jwt.AuthService, userRepository repositories.UserRepository) UserService {
	return &userService{
		redisApi:       redisApi,
		authService:    authService,
		userRepository: userRepository,
	}
}

func (s *userService) SetPassword(c context.Context, dto *dtos.SetPasswordDto) error {
	userBase, err := s.userRepository.FindUserBaseByUserId(c, dto.UserId)
	if err != nil {
		return err
	}
	if userBase.UserId == 0 {
		return xerror.NewXErrorByCode(err_codes.UserNotRegistered)
	}
	if dto.Password != "" || userBase.Password != "" {
		err = bcrypt.CompareHashAndPassword([]byte(userBase.Password), []byte(dto.Password))
		if err != nil {
			return xerror.NewXErrorByCode(err_codes.UserPasswordIncorrect)
		}
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userBase.Password = string(passwordByte)
	err = s.userRepository.UpdateUserBase(c, &userBase)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) LoginByPhoneCode(c context.Context, dto *dtos.LoginByPhoneCodeDto) (*dtos.UserLoginDto, error) {
	// 校验验证码
	code, err := s.redisApi.Get(c, consts.SEND_SMS_BASE_RDB_KER+dto.Phone, 0)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if dto.Code != code {
		// 这里可能换成code等
		return nil, xerror.NewXErrorByCode(err_codes.UserVerificationCodeIncorrect)
	}

	userLoginDto := &dtos.UserLoginDto{}

	// 获取用户信息
	userBase, err := s.userRepository.FindUserBaseByPhone(c, dto.Phone)
	if err != nil {
		return nil, err
	}
	// 用户不存在则创建(免了注册流程)
	if userBase.UserId == 0 {
		userBase, err = s.userRepository.CreateUserBaseByPhone(c, dto.Phone)
		if err != nil {
			return nil, err
		}
		userLoginDto.IsNew = true
	}
	userInfo, err := toUserInfoDto(&userBase)
	if err != nil {
		return nil, err
	}

	// 生成jwt_token
	jtwToken, err := s.generateJwtTokenByUserInfo(&userBase)
	if err != nil {
		return nil, err
	}
	userLoginDto.UserInfo = *userInfo
	userLoginDto.JwtToken = jtwToken
	return userLoginDto, nil
}

func (s *userService) LoginByPassword(c context.Context, dto *dtos.LoginByPassWordDto) (*dtos.UserLoginDto, error) {
	userBase, err := s.userRepository.FindUserBaseByPhone(c, dto.Phone)
	if err != nil {
		return nil, err
	}
	// 用户不存在则创建(免了注册流程)
	if userBase.UserId == 0 {
		return nil, xerror.NewXErrorByCode(err_codes.UserNotRegistered)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userBase.Password), []byte(dto.Password))
	if err != nil {
		return nil, xerror.NewXErrorByCode(err_codes.UserPasswordOrPhoneIncorrect)
	}

	userInfo, err := toUserInfoDto(&userBase)
	if err != nil {
		return nil, err
	}

	// 生成jwt_token
	jtwToken, err := s.generateJwtTokenByUserInfo(&userBase)
	if err != nil {
		return nil, err
	}

	userLoginDto := &dtos.UserLoginDto{
		UserInfo: *userInfo,
		JwtToken: jtwToken,
	}
	return userLoginDto, nil
}

func (s *userService) EditUserInfo(c context.Context, dto *dtos.UserInfoDto) error {
	/*
		text := fmt.Sprintf("%s , %s", dto.Username, dto.UserProfile)

		censorTextResult, err := s.censorGrpcServiceClient.CensorGrpcServiceClient.CensorText(c, &censorPb.CensorTextRequest{Text: text})
		if err != nil {
			return err
		}

		if censorTextResult.GetData().GetInterceptStatus() {
			return errors2.New(errors2.ErrParamInvalid).WithParam("存在敏感内容  " + s.censorGrpcServiceClient.AggregateCensorResult(censorTextResult))
		}

		censorImageResult, err := s.censorGrpcServiceClient.CensorGrpcServiceClient.CensorImage(c, &censorPb.CensorImageRequest{Image: dto.HeadPortrait})
		if err != nil {
			return err
		}

		if censorImageResult.GetData().GetInterceptStatus() {
			return errors2.New(errors2.ErrParamInvalid).WithParam("存在敏感内容  " + s.censorGrpcServiceClient.AggregateCensorResult(censorImageResult))
		}
	*/
	userInfo, err := fromUserInfoDto(dto)
	if err != nil {
		return err
	}
	return s.userRepository.UpdateUserBase(c, userInfo)
}

func (s *userService) GetUserInfoByUserId(c context.Context, req *dtos.UserIdDto) (*dtos.UserInfoDto, error) {
	keyPrefix := fmt.Sprintf("%s%d", consts.USER_INFO_RDB_KEY, req.UserId)
	userBase := mysql_models.TUserBase{}

	err := s.redisApi.GetAndUnmarshal(c, keyPrefix, userBase, 0)
	if err != nil || userBase.UserId == 0 {
		userBase, err = s.userRepository.FindUserBaseByUserId(c, req.UserId)
		if err != nil {
			return nil, err
		}
		_, err = s.redisApi.SetInterface(c, keyPrefix, userBase, consts.RDB_TIMEOUT, 0)
		if err != nil {
			zlog.Error("GetUserInfoByUserId cache.Set err : ", zap.Error(err))
		}
	}

	dto, err := toUserInfoDto(&userBase)
	dto.Phone = ""
	return dto, err
}

func (s *userService) LoginByOauthUserCode(c context.Context, dto *dtos.LoginByOauthUserCodeDto) (*dtos.UserLoginDto, error) {
	userLoginDto := &dtos.UserLoginDto{}

	// 获取用户信息
	userBase, err := s.userRepository.FindOauthUserBaseByCode(c, dto.Code)
	if err != nil {
		return nil, err
	}
	// 用户不存在则创建(免了注册流程)
	if userBase.UserId == 0 {
		userBase, err = s.userRepository.CreateOauthUserBaseByCode(c, dto.Code)
		if err != nil {
			return nil, err
		}
		userLoginDto.IsNew = true
	}
	userInfo, err := toUserInfoDto(&userBase)
	if err != nil {
		return nil, err
	}

	// 生成jwt_token
	jtwToken, err := s.generateJwtTokenByUserInfo(&userBase)
	if err != nil {
		return nil, err
	}

	userLoginDto.UserInfo = *userInfo
	userLoginDto.JwtToken = jtwToken
	return userLoginDto, nil
}

func (s *userService) generateJwtTokenByUserInfo(userBase *mysql_models.TUserBase) (string, error) {
	// 生成jwt_token
	claims := &jwt.Claims{
		JwtUserInfo: jwt.JwtUserInfo{
			UserId: userBase.UserId,
		},
	}
	return s.authService.CreateToken(claims)
}

func toUserInfoDto(userBase *mysql_models.TUserBase) (*dtos.UserInfoDto, error) {
	if userBase == nil {
		return nil, errors.New("userBase is nil")
	}
	dto := &dtos.UserInfoDto{
		UserId:       userBase.UserId,
		Phone:        userBase.Phone,
		Username:     userBase.Username,
		HeadPortrait: userBase.HeadPortrait,
		UserProfile:  userBase.UserProfile,
		Sex:          userBase.Sex,
	}
	return dto, nil
}

func fromUserInfoDto(userBase *dtos.UserInfoDto) (*mysql_models.TUserBase, error) {
	if userBase == nil {
		return nil, errors.New("userBase is nil")
	}
	dto := &mysql_models.TUserBase{
		UserId:       userBase.UserId,
		Username:     userBase.Username,
		HeadPortrait: userBase.HeadPortrait,
		UserProfile:  userBase.UserProfile,
		Sex:          userBase.Sex,
	}
	return dto, nil
}
