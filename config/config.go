package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Conf struct {
	Start_year int `yaml:"start_year"`
}

func (c *Conf) GetConf() * Conf{
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
