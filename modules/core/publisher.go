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

// Publisher is the one who wants to show an ad
type Publisher struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func publisherRoutes(router *gin.Engine) {
	publisher := router.Group("/api/publisher")
	{
		publisher.POST("/create", createPublisherHandler)
		publisher.GET("/", getAllPublishersHandler)
	}
}

func createPublisherHandler(c *gin.Context) {

	var publisher Publisher
	err := c.ShouldBindJSON(&publisher)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Missing required parameters",
		})
		return
	}

	keyUUID := uuid.New().String()
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "publishers", keyUUID)
	binName := as.NewBin("name", publisher.Name)
	binDescription := as.NewBin("description", publisher.Description)
	binActive := as.NewBin("isActive", 1)
	binUpdatedAt := as.NewBin("updatedAt", int32(time.Now().Unix()))
	binCreatedAt := as.NewBin("createdAt", int32(time.Now().Unix()))

	err = master.Client.PutBins(master.Client.DefaultWritePolicy, key, binName, binDescription, binActive, binUpdatedAt, binCreatedAt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "Publisher created",
			"key": keyUUID,
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
			result = append(result, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"publishers": result,
	})
}
