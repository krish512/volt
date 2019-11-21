package console

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes is route
func Routes(router *gin.Engine) {
	console := router.Group("/app")
	{
		console.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		})
	}
}
