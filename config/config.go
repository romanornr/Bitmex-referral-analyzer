package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"github.com/spf13/viper"
)

type Conf struct {
	Start_year int `yaml:"start_year"`
}

func (c *Conf) GetConf() *Conf {
	ymlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("file config.yml does not exist: %s\n", err)
	}

	err = yaml.Unmarshal(ymlFile, c)
	if err != nil {
		log.Fatalf("unable to unmarshal file: %s\n", err)
	}
	return c
}

// GetViperConfig reads config.yml file with viper
// example conf := config.GetViperConfig
// ip := viper.GetString("server.ip)
func GetViperConfig() error {
	viper.SetConfigType("yml")
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")


	err := viper.ReadInConfig()
	log.Println("Reading configuration from %s\n", viper.ConfigFileUsed())

	if err != nil {
		log.Fatalf("No configuration file loaded !\n%s", err)
	}
	return err
}