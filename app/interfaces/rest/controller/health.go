package controller

import (
	"github.com/gin-gonic/gin"
)

// HealthController ...
type HealthController struct{}

// ConsulHealthCheck ...
func (HealthController) ConsulHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"health": "nice",
	})

	return
}
