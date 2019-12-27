package master

import (
	"log"
	"time"
	"volt/config"

	"github.com/aerospike/aerospike-client-go"
	as "github.com/aerospike/aerospike-client-go"
)

// Client is connection struct for master database
var Client *as.Client

// ConnectMaster establishes connection with the master database
func ConnectMaster() {

	// port, _ := strconv.Atoi(config.Conf.Database.Aerospike.Port)
	client, err := as.NewClient(config.Conf.Database.Aerospike.Host, config.Conf.Database.Aerospike.Port)
	if err != nil {
		log.Fatalf("Aerospike connection failed: %v", err)
	}

	log.Printf("Connected to master database!")

	// Initialize default policies

	// Base Policy
	basePolicy := as.NewPolicy()
	basePolicy.TotalTimeout = config.Conf.Database.Aerospike.Timeout * time.Millisecond
	basePolicy.Priority = as.MEDIUM
	basePolicy.SendKey = true
	client.DefaultPolicy = basePolicy

	// Write Policy
	writePolicy := as.NewWritePolicy(0, 0)
	writePolicy.TotalTimeout = config.Conf.Database.Aerospike.Timeout * time.Millisecond
	writePolicy.Priority = as.HIGH
	writePolicy.SendKey = true
	writePolicy.RecordExistsAction = as.UPDATE
	client.DefaultWritePolicy = writePolicy

	// Query Policy
	queryPolicy := as.NewQueryPolicy()
	queryPolicy.TotalTimeout = config.Conf.Database.Aerospike.Timeout * time.Millisecond
	queryPolicy.Priority = as.LOW
	queryPolicy.SendKey = true
	queryPolicy.IncludeBinData = true
	client.DefaultQueryPolicy = queryPolicy

	// Test database
	key, _ := as.NewKey(config.Conf.Database.Aerospike.Namespace, "core", "meta")
	bin := as.NewBin("timestamp", int32(time.Now().Unix()))

	err = client.PutBins(writePolicy, key, bin)
	if err != nil {
		log.Fatalf("Aerospike write error: %v", err)
	}

	Client = client
	createIndexes()
}

func createIndexes() {
	Client.CreateIndex(Client.DefaultWritePolicy, config.Conf.Database.Aerospike.Namespace, "advertisers", "idx_advertisers_isactive", "isActive", aerospike.NUMERIC)
	Client.CreateIndex(Client.DefaultWritePolicy, config.Conf.Database.Aerospike.Namespace, "campaigns", "idx_campaigns_isactive", "isActive", aerospike.NUMERIC)
	Client.CreateIndex(Client.DefaultWritePolicy, config.Conf.Database.Aerospike.Namespace, "campaigns", "idx_campaigns_status", "status", aerospike.STRING)
	Client.CreateIndex(Client.DefaultWritePolicy, config.Conf.Database.Aerospike.Namespace, "adunits", "idx_adunits_isactive", "isActive", aerospike.NUMERIC)
	Client.CreateIndex(Client.DefaultWritePolicy, config.Conf.Database.Aerospike.Namespace, "publishers", "idx_publishers_isactive", "isActive", aerospike.NUMERIC)
}
