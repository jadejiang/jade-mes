package app

import (
	"time"

	"github.com/gin-gonic/gin"

	"jade-mes/app/infrastructure/env"
	"jade-mes/config"
	"jade-mes/consul"
)

var router *gin.Engine
var port string

func init() {
	println("initing app...")
	settings := config.GetConfig()

	env.StartTime = time.Now()

	isReleaseMode := settings.GetBool("release_mode")
	if isReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router = NewRouter()
	port = settings.GetString("server.port")
}

// Start starts gin app
func Start() {
	if config.Config.Consul.Enable {
		consul.RegisterServer()
	}
	router.Run(port)
}
