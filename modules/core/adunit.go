package core

import (
	"log"
	"net/http"
	"time"
	"volt/config"
	"volt/database/master"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func adunitRoutes(router *gin.Engine) {
	adunit := router.Group("/api/adunit")
	{
		adunit.POST("/create", createAdUnitHandler)
		adunit.GET("/", getAllAdUnitsHandler)
	}
}

func createAdUnitHandler(c *gin.Context) {

	keyUUID := uuid.New().String()
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "adunits", keyUUID)
	binName := as.NewBin("name", "Paragon")
	binDescription := as.NewBin("description", c.Query("description"))
	binPublisher := as.NewBin("publisher", "54ea76a8-1537-11ea-8cd4-a683e78ecd24")
	binActive := as.NewBin("isActive", 1)
	binUpdatedAt := as.NewBin("updatedAt", int32(time.Now().Unix()))
	binCreatedAt := as.NewBin("createdAt", int32(time.Now().Unix()))

	err := master.Client.PutBins(master.Client.DefaultWritePolicy, key, binName, binDescription, binPublisher, binActive, binUpdatedAt, binCreatedAt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "Ad Unit created",
		})
	}
}

func getAllAdUnitsHandler(c *gin.Context) {

	stmt := as.NewStatement(config.Conf.Database.Aerospike.Namespace, "adunits", "name", "description")
	stmt.SetFilter(as.NewEqualFilter("isActive", 1))
	rs, err := master.Client.Query(master.Client.DefaultQueryPolicy, stmt)
	if err != nil {
		log.Fatalf("Master query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	}

	var result []map[string]string
	for res := range rs.Results() {
		if res.Err != nil {
			log.Printf("Error: %v", err)
		} else {
			item := make(map[string]string)
			item["key"] = res.Record.Key.Value().String()
			item["name"] = res.Record.Bins["name"].(string)
			item["description"] = res.Record.Bins["description"].(string)
			result = append(result, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"adunits": result,
	})
}
