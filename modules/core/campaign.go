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

func campaignRoutes(router *gin.Engine) {
	campaign := router.Group("/api/campaign")
	{
		campaign.POST("/create", createCampaignHandler)
		campaign.GET("/", getAllCampaignsHandler)
	}
}

func createCampaignHandler(c *gin.Context) {

	keyUUID := uuid.Must(uuid.NewUUID()).String()
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "campaigns", keyUUID)
	binName := as.NewBin("name", "Paragon")
	binDescription := as.NewBin("description", "Our premium campaign")
	binStatus := as.NewBin("status", "scheduled")
	binActive := as.NewBin("isActive", 1)
	binStartAt := as.NewBin("startAt", int32(time.Now().Unix()))
	binEndAt := as.NewBin("endAt", int32(time.Now().Unix()))
	binAdvertiser := as.NewBin("advertiser", "54ea76a8-1537-11ea-8cd4-a683e78ecd24")
	binUpdatedAt := as.NewBin("updatedAt", int32(time.Now().Unix()))
	binCreatedAt := as.NewBin("createdAt", int32(time.Now().Unix()))

	err := master.Client.PutBins(master.Client.DefaultWritePolicy, key, binName, binDescription, binStatus, binActive, binStartAt, binEndAt, binAdvertiser, binUpdatedAt, binCreatedAt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "Create Campaign",
		})
	}
}

func getAllCampaignsHandler(c *gin.Context) {

	stmt := as.NewStatement(config.Conf.Database.Aerospike.Namespace, "campaigns", "name", "description", "status")
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
		"campaigns": result,
	})
}
