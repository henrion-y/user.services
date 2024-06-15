package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/henrion-y/base.services/http/gin/middlewares"
	"go.uber.org/dig"

	"user.services/app/http/controllers"
	middlewares2 "user.services/app/http/middlewares"
)

func APIRoutes(c RouterContext) {
	api := c.Router.Group("/api/v1/users")
	{
		sms := api.Group("/sms")
		{
			sms.POST("/send_code", c.LimitMiddleware.Handler(1, 1, 3*time.Second), c.SmsController.SendCode) // 发送验证码
		}
		user := api.Group("/user")
		{
			user.POST("/login_by_pass_word", c.UserController.LoginByPassWord)                                                 // 账号密码登录
			user.POST("/login_by_phone_code", c.UserController.LoginByPhoneCode)                                               // 验证码登录
			user.POST("/login_by_oauth_code", c.UserController.LoginByOauthUserCode)                                           // 第三方code登录
			user.GET("/get_user_info_by_user_id", c.UserController.GetUserInfoByUserId)                                        // 根据用户id获取用户信息
			user.POST("/set_password", c.AuthMiddleware.SetClaimsAbortTourist(), c.UserController.SetPassword)                 // 重置密码
			user.POST("/edit_user_info", c.AuthMiddleware.SetClaimsAbortTourist(), c.UserController.EditUserInfo)              // 修改用户信息
			user.GET("/get_user_info_by_token", c.AuthMiddleware.SetClaimsAbortTourist(), c.UserController.GetUserInfoByToken) // 修改用户信息

		}
	}
}

type RouterContext struct {
	dig.In
	Router          *gin.Engine
	LimitMiddleware *middlewares2.LimitMiddleware
	AuthMiddleware  *middlewares.JWTAuthMiddleware
	SmsController   *controllers.SmsController
	UserController  *controllers.UserController
}
