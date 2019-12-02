package config

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

// Conf is struct to parse configuration file
type conf struct {
	Database struct {
		Aerospike struct {
			Host      string        `yaml:"host"`
			Port      int           `yaml:"port"`
			Timeout   time.Duration `yaml:"timeout"`
			Namespace string        `yaml:"namespace"`
		}
	}
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Conf has configurations from the yaml file at default path
var Conf conf

// InitConf reads configuration from the yaml file at default path and initialised Conf variable
func InitConf() {

	var c conf

	config, err := ioutil.ReadFile("config/volt.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(config, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	Conf = c
}
