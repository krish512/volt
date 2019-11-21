package main

import (
	"fmt"
	"net/http"

	"volt/config"
	"volt/modules/console"
	"volt/modules/core"
	"volt/modules/delivery"

	"github.com/gin-gonic/gin"
)

func commonHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("server", "Volt")
	}
}

func main() {

	var conf config.Conf
	fmt.Println(conf.GetConf())

	router := gin.Default()
	router.Use(commonHeaders())
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "online",
		})
	})

	console.Routes(router)
	core.Routes(router)
	delivery.Routes(router)

	router.Run(":" + conf.GetConf().Port)
}
