package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes has all the core routes
func Routes(router *gin.Engine) {
	core := router.Group("/api")
	{
		core.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"module": "core",
				"msg":    "online",
			})
		})
	}

	// Append Advertiser routes
	advertiserRoutes(router)
	campaignRoutes(router)
	inventoryRoutes(router)
	publisherRoutes(router)
}
