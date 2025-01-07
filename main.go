// @title           Learning Golang
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"learn/api"
	redis "learn/db/redis"
	"learn/settings"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	if settings.Settings.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = api.ApiInit()
	redis.MakeClient()
}

func main() {
	router.Run(":8015")
}
