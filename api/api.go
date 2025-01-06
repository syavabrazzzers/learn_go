package api

import (
	"learn/api/endpoints"
	"learn/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"
)

func ApiInit() *gin.Engine {
	api := gin.New()

	api.Use(gin.Logger())

	api.SetTrustedProxies([]string{"localhost"})

	docs.SwaggerInfo.BasePath = "/api"

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	endpoints.AddUsers(api)
	endpoints.AddAuth(api)

	return api
}
