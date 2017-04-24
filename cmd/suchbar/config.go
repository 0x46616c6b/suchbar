package main

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Page contains the ID and an Alias
type Page struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}

type Grafana struct {
	Host        string `yaml:"host"`
	APIKey      string `yaml:"api_key"`
	ElasticHost string `yaml:"elastic_host"`
}

type Elastic struct {
	Host string `yaml:"host"`
}

//Config contains all Configuration values for suchbar
type Config struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	Elastic   Elastic `yaml:"elastic"`
	Grafana   Grafana `yaml:"grafana"`
	LogLevel  string `yaml:"log_level"`
	Pages     []Page `yaml:"pages"`
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
