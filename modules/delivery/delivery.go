package delivery

import (
	"log"
	"net/http"
	"runtime/debug"
	"volt/config"
	"volt/database/master"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/gin-gonic/gin"
)

// CampaignDeliver is structure to store campaign details for delivery
type CampaignDeliver struct {
	Campaign  string        `json:"campaign" binding:"required"`
	Creatives []interface{} `json:"creatives" binding:"required"`
}

var campaignsCache map[string]CampaignDeliver

// Routes has all the delivery routes
func Routes(router *gin.Engine) {
	serve := router.Group("/serve")
	{
		serve.GET("/ad/:adunit/call", getAdHandler)
		serve.GET("/reg/:campaign/:event/ack.png", trackEventHandler)
	}
}

// UpdateActiveCampaignsCache regularly updates the campaignsCache object for delivery
func UpdateActiveCampaignsCache() {
	// Update status to active for all scheduled campaigns where startAt < now
	// Update status to complete for all active campaigns where endAt < now
	// Select all campaigns where status = active

	// Update active campaigns cache with the new object

	// ToDo

	// PoC
	var campaignsCacheTmp map[string]CampaignDeliver
	campaignsCacheTmp = make(map[string]CampaignDeliver)
	stmt := as.NewStatement(config.Conf.Database.Aerospike.Namespace, "campaigns", "adunits", "creatives")
	stmt.SetFilter(as.NewEqualFilter("isActive", 1))
	// stmt.SetFilter(as.NewEqualFilter("startAt", 1))
	// stmt.SetFilter(as.NewEqualFilter("status", 1))
	rs, err := master.Client.Query(master.Client.DefaultQueryPolicy, stmt)
	if err != nil {
		log.Fatalf("Master query error: %v", err)
	}

	// var result []map[string]interface{}
	for res := range rs.Results() {
		if res.Err != nil {
			log.Printf("Error: %v", err)
		} else {
			campaignKey := res.Record.Key.Value().String()
			creatives := res.Record.Bins["creatives"].([]interface{})
			campaignDeliver := CampaignDeliver{
				Campaign:  campaignKey,
				Creatives: creatives,
			}
			adunits := res.Record.Bins["adunits"].([]interface{})
			for _, adunit := range adunits {
				campaignsCacheTmp[adunit.(string)] = campaignDeliver
			}
		}
	}

	// Update campaignsCache
	campaignsCache = campaignsCacheTmp
	campaignsCacheTmp = nil

	// Finally, free up memory
	debug.FreeOSMemory()
}

func getAdHandler(c *gin.Context) {
	adunit := c.Param("adunit")
	c.JSON(http.StatusOK, gin.H{
		"result": campaignsCache[adunit],
	})
}

func trackEventHandler(c *gin.Context) {
	campaign := c.Param("campaign")
	event := c.Param("event")
	c.JSON(http.StatusOK, gin.H{
		"msg":      "ok",
		"campaign": campaign,
		"event":    event,
	})
}
