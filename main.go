package main

import (
	"github.com/chuongthanh0410/interview/app"
)

// @title Gin Swagger Example API
// @version 1.0
// @description API xác thực dùng token tĩnh, giá trị = 1234567890abcdefjustforspeed
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @host localhost:8080
// @BasePath /api/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @securityDefinitions.jwt BearerAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
