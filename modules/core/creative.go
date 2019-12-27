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

type creative struct {
	name        string
	description string
	kind        string
	url         string
}

func creativeRoutes(router *gin.Engine) {
	creative := router.Group("/api/creative")
	{
		creative.POST("/create", createCreativeHandler)
		creative.GET("/", getAllCreativesHandler)
	}
}

func createCreativeHandler(c *gin.Context) {

	keyUUID := uuid.New().String()
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "creatives", keyUUID)
	binName := as.NewBin("name", "Paragon")
	binDescription := as.NewBin("description", "Our premium creative")
	binKind := as.NewBin("kind", "URL")
	binURL := as.NewBin("url", "https://encrypted-tbn0.gstatic.com/images?q=tbn%3AANd9GcSqff08PL_wxf7-F5WkuaK5w3DQf5aWbW4jD4MyZvd6w4kkFQim")
	binUpdatedAt := as.NewBin("updatedAt", int32(time.Now().Unix()))
	binCreatedAt := as.NewBin("createdAt", int32(time.Now().Unix()))

	err := master.Client.PutBins(master.Client.DefaultWritePolicy, key, binName, binDescription, binKind, binURL, binUpdatedAt, binCreatedAt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "Creative created",
			"key": keyUUID,
		})
	}
}

func getAllCreativesHandler(c *gin.Context) {

	stmt := as.NewStatement(config.Conf.Database.Aerospike.Namespace, "creatives", "name", "description", "kind", "url")
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
			for key := range res.Record.Bins {
				item[key] = res.Record.Bins[key].(string)
			}
			result = append(result, item)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"creatives": result,
	})
}
