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

func publisherRoutes(router *gin.Engine) {
	publisher := router.Group("/api/publisher")
	{
		publisher.POST("/create", createPublisherHandler)
		publisher.GET("/", getAllPublishersHandler)
	}
}

func createPublisherHandler(c *gin.Context) {

	keyUUID := uuid.Must(uuid.NewUUID()).String()
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "publishers", keyUUID)
	binName := as.NewBin("name", "Paragon")
	binDescription := as.NewBin("description", "Our premium publisher")
	binActive := as.NewBin("isActive", 1)
	binUpdatedAt := as.NewBin("updatedAt", int32(time.Now().Unix()))
	binCreatedAt := as.NewBin("createdAt", int32(time.Now().Unix()))

	err := master.Client.PutBins(master.Client.DefaultWritePolicy, key, binName, binDescription, binActive, binUpdatedAt, binCreatedAt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "Create publisher",
		})
	}
}

func getAllPublishersHandler(c *gin.Context) {

	stmt := as.NewStatement(config.Conf.Database.Aerospike.Namespace, "publishers", "name", "description", "status")
	stmt.SetFilter(as.NewEqualFilter("isActive", 1))
	rs, err := master.Client.Query(master.Client.DefaultQueryPolicy, stmt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
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
			item["status"] = res.Record.Bins["status"].(string)
			result = append(result, item)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"publishers": result,
	})
}
