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

// Campaign is used to serve and track ad creatives
type Campaign struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	StartAt     int32    `json:"startAt" binding:"required" time_format:"unix"`
	EndAt       int32    `json:"endAt" time_format:"unix"`
	Advertiser  string   `json:"advertiser" binding:"required"`
	Creatives   []string `json:"creatives"`
}

func campaignRoutes(router *gin.Engine) {
	campaign := router.Group("/api/campaign")
	{
		campaign.POST("/assign/:campaign/advertiser", assignAdvertiserCampaignHandler)
		campaign.POST("/assign/:campaign", assignAdUnitCampaignHandler)
		campaign.POST("/create", createCampaignHandler)
		campaign.GET("/", getAllCampaignsHandler)
	}
}

func createCampaignHandler(c *gin.Context) {

	var campaign Campaign
	err := c.ShouldBindJSON(&campaign)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Missing required parameters",
		})
		return
	}

	keyUUID := uuid.New().String()
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "campaigns", keyUUID)
	binName := as.NewBin("name", campaign.Name)
	binDescription := as.NewBin("description", campaign.Description)
	binStatus := as.NewBin("status", "inactive")
	binActive := as.NewBin("isActive", 1)
	binStartAt := as.NewBin("startAt", campaign.StartAt)
	binEndAt := as.NewBin("endAt", campaign.EndAt)
	binAdvertiser := as.NewBin("advertiser", campaign.Advertiser)
	binCreatives := as.NewBin("creatives", []creative{})
	binAdunits := as.NewBin("adunits", []string{})
	binUpdatedAt := as.NewBin("updatedAt", int32(time.Now().Unix()))
	binCreatedAt := as.NewBin("createdAt", int32(time.Now().Unix()))

	err = master.Client.PutBins(master.Client.DefaultWritePolicy, key, binName, binDescription, binStatus, binActive, binStartAt, binEndAt, binAdvertiser, binCreatives, binAdunits, binUpdatedAt, binCreatedAt)
	if err != nil {
		log.Fatalf("Master write error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"msg": "Campaign created",
			"key": keyUUID,
		})
	}
}

// assignAdUnitCampaignHandler adds adunit to the specified campaign
func assignAdUnitCampaignHandler(c *gin.Context) {

	// Parse campaign id from url
	campaign := c.Param("campaign")

	// Get AdUnit id from request body
	type Adunit struct {
		Adunit string `json:"adunit" binding:"required"`
	}
	var adunit Adunit
	bindErr := c.ShouldBindJSON(&adunit)
	if bindErr != nil {
		log.Println(bindErr)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Missing required parameters",
		})
		return
	}

	// Check if campaign exists
	campaignKey, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "campaigns", campaign)
	campaignExists, campaignErr := master.Client.Exists(master.Client.DefaultPolicy, campaignKey)
	if campaignErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   campaignErr.Error(),
		})
		return
	}
	if !campaignExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Campaign does not exist",
		})
		return
	}

	// // Check if adunit exists
	// adunitKey, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "adunits", adunit.Adunit)
	// adunitExists, adunitErr := master.Client.Exists(master.Client.DefaultPolicy, adunitKey)
	// if adunitErr != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": true,
	// 		"msg":   adunitErr.Error(),
	// 	})
	// 	return
	// }
	// if !adunitExists {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": true,
	// 		"msg":   "Ad Unit does not exist",
	// 	})
	// 	return
	// }

	// Add adunit to the campaign
	listPolicy := as.NewListPolicy(as.ListOrderOrdered, as.ListWriteFlagsAddUnique)
	appendAdUnit := as.ListAppendWithPolicyOp(listPolicy, "adunits", adunit.Adunit)
	records, err := master.Client.Operate(master.Client.DefaultWritePolicy, campaignKey, appendAdUnit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg": records.Bins["adunits"].(int),
	})
}

// assignAdvertiserCampaignHandler adds advertiser to the specified campaign
func assignAdvertiserCampaignHandler(c *gin.Context) {

	// Parse campaign id from url
	campaign := c.Param("campaign")

	// Get Advertiser id from request body
	type Advertiser struct {
		Advertiser string `json:"advertiser" binding:"required"`
	}
	var advertiser Advertiser
	bindErr := c.ShouldBindJSON(&advertiser)
	if bindErr != nil {
		log.Println(bindErr)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Missing required parameters",
		})
		return
	}

	// Check if campaign exists
	campaignKey, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "campaigns", campaign)
	campaignExists, campaignErr := master.Client.Exists(master.Client.DefaultPolicy, campaignKey)
	if campaignErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   campaignErr.Error(),
		})
		return
	}
	if !campaignExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Campaign does not exist",
		})
		return
	}

	// Check if advertiser exists
	advertiserKey, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "advertisers", advertiser.Advertiser)
	advertiserExists, advertiserErr := master.Client.Exists(master.Client.DefaultPolicy, advertiserKey)
	if advertiserErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   advertiserErr.Error(),
		})
		return
	}
	if !advertiserExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Advertiser does not exist",
		})
		return
	}

	// Add advertiser to the campaign
	binAdvertiser := as.NewBin("advertiser", advertiser.Advertiser)
	err := master.Client.PutBins(master.Client.DefaultWritePolicy, campaignKey, binAdvertiser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":        "Advertiser assigned",
		"advertiser": advertiser.Advertiser,
		"campaign":   campaign,
	})
}

func getAllCampaignsHandler(c *gin.Context) {

	stmt := as.NewStatement(config.Conf.Database.Aerospike.Namespace, "campaigns", "name", "description", "advertiser", "adunits", "status", "startAt", "endAt")
	stmt.SetFilter(as.NewEqualFilter("isActive", 1))
	rs, err := master.Client.Query(master.Client.DefaultQueryPolicy, stmt)
	if err != nil {
		log.Fatalf("Master query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err,
		})
	}

	var result []map[string]interface{}
	for res := range rs.Results() {
		if res.Err != nil {
			log.Printf("Error: %v", err)
		} else {
			item := make(map[string]interface{})
			item["key"] = res.Record.Key.Value().String()
			item["name"] = res.Record.Bins["name"].(string)
			item["description"] = res.Record.Bins["description"].(string)
			item["advertiser"] = res.Record.Bins["advertiser"].(string)
			item["adunits"] = res.Record.Bins["adunits"].([]interface{})
			item["status"] = res.Record.Bins["status"].(string)
			item["startAt"] = res.Record.Bins["startAt"].(int)
			item["endAt"] = res.Record.Bins["endAt"].(int)
			result = append(result, item)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"campaigns": result,
	})
}
