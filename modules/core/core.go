package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes is route
func Routes(router *gin.Engine) {
	core := router.Group("/core")
	{
		core.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		})
	}
}
