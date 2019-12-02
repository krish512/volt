package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"volt/config"
	"volt/database/master"
	"volt/modules/console"
	"volt/modules/core"
	"volt/modules/delivery"
	"volt/utils"

	"github.com/gin-gonic/gin"
)

func commonHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("server", "Volt")
	}
}

func initResources() {
	defer utils.Logger.Sync()
	config.InitConf()
	master.ConnectMaster()
}

func closeResources() {
	defer utils.Logger.Sync()
	master.Client.Close()
	println("Stopping server!")
}

func main() {

	// Handling SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		closeResources()
		os.Exit(1)
	}()

	// Initialise
	initResources()

	// Declare routes
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(commonHeaders())
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "online",
		})
	})

	// Load rotues from modules
	console.Routes(router)
	core.Routes(router)
	delivery.Routes(router)

	// Handle 404 routes
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
	})

	// Start server
	log.Printf("Starting server at " + ":" + config.Conf.Port)
	router.Run(":" + config.Conf.Port)
}
