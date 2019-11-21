package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

// Routes is route
func Routes(router *gin.Engine) {
	serve := router.Group("/serve")
	{
		serve.GET("/ad/:target/call", dummy)
		serve.GET("/reg/:campaign/:event/ack.png", dummy)
	}
}
