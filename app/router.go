package app

import (
	"jade-mes/config"

	"github.com/gin-contrib/expvar"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"jade-mes/app/infrastructure/middleware"
	"jade-mes/app/interfaces/rest/controller"
)

const (
	// DefaultMetricPath ...
	DefaultMetricPath = "/metrics/"
	// DefaultDebugPath ...
	DefaultDebugPath = "/thisisdebug"
	// HealthPath ...
	HealthPath = "/health"
)

// NewRouter initializes router with controllers and middlewares
func NewRouter() *gin.Engine {
	router := gin.New()

	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	router.Use(middleware.Logger(DefaultMetricPath, HealthPath))
	router.Use(middleware.Prometheus(DefaultMetricPath))
	router.Use(middleware.TraceMiddleware(config.Config.Tracer.Name))
	// middleware.InitMetrics(router)

	metrics := router.Group(DefaultMetricPath)
	{
		metrics.GET("/heath", middleware.MetricHandler)
	}

	c := new(controller.HealthController)

	router.GET(HealthPath, c.ConsulHealthCheck)

	debug := router.Group(DefaultDebugPath)
	{
		debug.GET("/vars", expvar.Handler())
		pprof.Register(router, DefaultDebugPath)
	}

	user := router.Group("mes/user/v1")
	{
		useController := new(controller.UserController)
		user.POST("/register", useController.Register)
	}

	role := router.Group("mes/role/v1")
	{
		roleController := new(controller.RoleController)
		role.POST("/createRole", roleController.CreateRole)
		role.POST("/detailById", roleController.FindRoleByID)
		role.POST("/detailByName", roleController.FindRoleByName)
	}

	return router
}
