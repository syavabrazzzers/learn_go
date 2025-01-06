package main

import (
	"learn/api"
	"learn/db"
	redis "learn/db/redis"
	"learn/settings"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	settings.Init()
	if settings.Settings.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = api.ApiInit()
	db.Init()
	redis.MakeClient()
}

func main() {
	router.Run(":8015")
}
