package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Conf is struct
type Conf struct {
	Database struct {
		Influx struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		}
	}
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// GetConf is get config
func (c *Conf) GetConf() *Conf {

	config, err := ioutil.ReadFile("config/volt.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(config, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
