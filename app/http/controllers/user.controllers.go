package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/henrion-y/base.services/infra/jwt"
	"user.services/app/domain/dtos"
	"user.services/app/domain/services"
)

type UserController struct {
	authService jwt.AuthService
	userService services.UserService
}

func NewUserController(authService jwt.AuthService, userService services.UserService) *UserController {
	return &UserController{
		authService: authService,
		userService: userService,
	}
}

// LoginByPhoneCode
// @Tags 用户
// @Summary 手机验证码登录
// @Description 手机验证码登录
// @Id LoginByPhoneCode
// @Produce  json
// @Param        LoginByPhoneCodeDto  body  dtos.LoginByPhoneCodeDto  true  "手机验证码登录"
// @Success 200 {object} ResponseData{data=dtos.UserLoginDto}
// @Router /user/login_by_phone_code [post]
func (h *UserController) LoginByPhoneCode(c *gin.Context) {
	req := &dtos.LoginByPhoneCodeDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}

	useLoginInfo, err := h.userService.LoginByPhoneCode(c, req)
	if err != nil {
		if err.Error() == "验证码错误" {
			responseError(c, err)
			return
		}
		responseError(c, err)
		return
	}
	responseData(c, useLoginInfo)
}

// LoginByOauthUserCode
// @Tags 用户
// @Summary 第三方code登录
// @Description 第三方code登录
// @Id LoginByOauthUserCode
// @Produce  json
// @Param        LoginByOauthUserCodeDto  body  dtos.LoginByOauthUserCodeDto  true  "第三方code登录"
// @Success 200 {object} ResponseData{data=dtos.UserLoginDto}
// @Router /user/login_by_oauth_code [post]
func (h *UserController) LoginByOauthUserCode(c *gin.Context) {
	req := &dtos.LoginByOauthUserCodeDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}

	useLoginInfo, err := h.userService.LoginByOauthUserCode(c, req)
	if err != nil {
		responseError(c, err)
		return
	}
	responseData(c, useLoginInfo)
}

// LoginByPassWord
// @Tags 用户
// @Summary 账号密码登录
// @Description 账号密码登录
// @Id LoginByPPassWord
// @Produce  json
// @Param        LoginByPassWordDto  body  dtos.LoginByPassWordDto  true  "账号密码登录"
// @Success 200 {object} ResponseData{data=dtos.UserLoginDto}
// @Router /user/login_by_pass_word [post]
func (h *UserController) LoginByPassWord(c *gin.Context) {
	req := &dtos.LoginByPassWordDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}

	useLoginInfo, err := h.userService.LoginByPassword(c, req)
	if err != nil {
		if err.Error() == "账号未注册" || err.Error() == "手机号或密码错误" {
			responseError(c, err)
			return
		}
		responseError(c, err)
		return
	}
	responseData(c, useLoginInfo)
}

// EditUserInfo
// @Tags 用户
// @Summary 编辑用户信息
// @Description 编辑用户信息
// @securityDefinitions.basic  BasicAuth
// @Id EditUserInfo
// @Produce  json
// @Param        UserInfoDto  body  dtos.UserInfoDto  true  "编辑用户信息"
// @Success 200 {object} ResponseData{}
// @Router /user/edit_user_info [post]
func (h *UserController) EditUserInfo(c *gin.Context) {
	claims, err := h.authService.GetClaimsByGinCtx(c)
	if err != nil {
		responseError(c, err)
		return
	}

	req := &dtos.UserInfoDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}
	req.UserId = claims.UserId

	err = h.userService.EditUserInfo(c, req)
	if err != nil {
		responseError(c, err)
		return
	}
	responseSuccess(c)
}

// SetPassword
// @Tags 用户
// @Summary 重置密码
// @Description 重置密码
// @securityDefinitions.basic  BasicAuth
// @Id SetPassword
// @Produce  json
// @Param        SetPasswordDto  body  dtos.SetPasswordDto  true  "重置密码"
// @Success 200 {object} ResponseData{}
// @Router /user/set_password [post]
func (h *UserController) SetPassword(c *gin.Context) {
	claims, err := h.authService.GetClaimsByGinCtx(c)
	if err != nil {
		responseError(c, err)
		return
	}

	req := &dtos.SetPasswordDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}

	req.UserId = claims.UserId
	err = h.userService.SetPassword(c, req)
	if err != nil {
		responseError(c, err)
		return
	}
	responseSuccess(c)
}

// GetUserInfoByUserId
// @Tags 用户
// @Summary 根据用户id获取用户信息
// @Description 根据用户id获取用户信息
// @securityDefinitions.basic  BasicAuth
// @Id GetUserInfoByUserId
// @Produce  json
// @Param        UserIdDto  query  dtos.UserIdDto  true  "根据用户id获取用户信息"
// @Success 200 {object} ResponseData{data=dtos.UserInfoDto}
// @Router /user/get_user_info_by_user_id [get]
func (h *UserController) GetUserInfoByUserId(c *gin.Context) {
	req := &dtos.UserIdDto{}
	if err := c.Bind(req); err != nil {
		responseError(c, err)
		return
	}

	userInfo, err := h.userService.GetUserInfoByUserId(c, req)
	if err != nil {
		responseError(c, err)
		return
	}
	responseData(c, userInfo)
}

// GetUserInfoByToken
// @Tags 用户
// @Summary 根据token获取用户信息
// @Description 根据token获取用户信息
// @securityDefinitions.basic  BasicAuth
// @Id GetUserInfoByToken
// @Produce  json
// @Success 200 {object} ResponseData{data=dtos.UserInfoDto}
// @Router /user/get_user_info_by_token [get]
func (h *UserController) GetUserInfoByToken(c *gin.Context) {
	claims, err := h.authService.GetClaimsByGinCtx(c)
	if err != nil {
		responseError(c, err)
		return
	}
	req := &dtos.UserIdDto{UserId: claims.UserId}
	userInfo, err := h.userService.GetUserInfoByUserId(c, req)
	if err != nil {
		responseError(c, err)
		return
	}
	responseData(c, userInfo)
}
