package bootstrap

import (
	"github.com/henrion-y/base.services/database/gorm"
	"github.com/henrion-y/base.services/domain/repository/gormrepo"
	middlewares2 "github.com/henrion-y/base.services/http/gin/middlewares"
	"github.com/henrion-y/base.services/infra/jwt"
	"github.com/henrion-y/base.services/infra/qiniu"
	"github.com/henrion-y/base.services/infra/redisapi"
	"go.uber.org/fx"
	"user.services/app/domain/repositories"
	"user.services/app/domain/services"
	"user.services/app/http/controllers"
	"user.services/app/http/middlewares"
	"user.services/app/providers"
	"user.services/configs"
	"user.services/pkg/sdks/shansuma/sms"
	"user.services/routes"
)

func App() *fx.App {
	options := make([]fx.Option, 0)

	comOptions := []fx.Option{
		// Configurations (./config)
		fx.Provide(configs.NewConfigProvider),

		// Providers (./app/providers)
		fx.Provide(providers.NewGinProvider),
		fx.Provide(providers.NewLoggerProvider),

		fx.Provide(jwt.NewJWTService),
		fx.Provide(sms.NewSmsProvider),
		fx.Provide(qiniu.NewQiNiuClient),

		// Providers (db)
		fx.Provide(redisapi.NewRedisApiProvider),
		fx.Provide(gorm.NewDbProvider),

		// Middlewares (./app/http/middlewares)
		fx.Provide(middlewares.NewLogMiddleware),
		fx.Provide(middlewares2.NewJWTAuthMiddleware),
		fx.Provide(middlewares.NewLimitMiddleware),

		// NerRepository
		// Repositories (./app/domain/repositories)
		fx.Provide(gormrepo.NewBaseRepository),
		fx.Provide(repositories.NewUserRepository),

		// NewService
		// Services (./app/domain/services)
		fx.Provide(services.NewUserService),
		fx.Provide(services.NewSmsService),

		// NewController
		// Controllers (./app/http/controllers)
		fx.Provide(controllers.NewUserController),
		fx.Provide(controllers.NewSmsController),

		fx.Invoke(middlewares.GlobalMiddlewares),
		fx.Invoke(routes.APIRoutes),
		fx.Invoke(providers.StartService),
	}

	options = append(options, comOptions...)

	return fx.New(options...)
}
