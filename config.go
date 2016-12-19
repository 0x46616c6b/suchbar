package main

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Page contains the ID and an Alias
type Page struct {
	ID    string `yaml:"id"`
	Alias string `yaml:"alias"`
}

//Config contains all Configuration values for suchbar
type Config struct {
	AppID       string `yaml:"app_id"`
	AppSecret   string `yaml:"app_secret"`
	ElasticHost string `yaml:"elastic_host"`
	LogLevel    string `yaml:"log_level"`
	Pages       []Page `yaml:"pages"`
}

func (c *Config) read(path string) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Cannot read config file")
	}

	err = yaml.Unmarshal(f, c)
	if err != nil {
		log.Fatal("Cannot parse config file")
	}
}
