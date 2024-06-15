package main

import (
	"context"
	"log"

	"user.services/bootstrap"
)

// @title user.services
// @version 2.0
// @description user.services API接口文档

// @contact.name API Support
// @contact.url http://www.swagger.io/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8060
// @BasePath /api/v1/users
// @schemes http
func main() {
	app := bootstrap.App()
	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
