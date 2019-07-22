package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"text/template"

	log "github.com/Sirupsen/logrus"

	"github.com/0x46616c6b/suchbar/grafana"
)

var (
	configPath   string
	templatePath string
	config       Config
)

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "path to the configuration file")
	flag.StringVar(&templatePath, "template", "dashboard.tpl", "path to dashboard template file")
	flag.Parse()

	config.read(configPath)

	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}

func main() {
	client, err := grafana.NewClient(config.Grafana.APIKey, config.Grafana.Host)
	if err != nil {
		log.Println(err)
	}

	log.Println("Create Datasources and dashboards")
	for _, page := range config.Pages {
		createDatasources(client, page)
		createDashboard(client, page)
	}
}

func createDatasources(client *grafana.Client, page Page) {
	options := map[string]interface{}{
		"type":      "elasticsearch",
		"name":      page.Name,
		"url":       config.Grafana.ElasticHost,
		"database":  page.Alias,
		"access":    "proxy",
		"basicAuth": false,
		"isDefault": false,
		"jsonData": map[string]interface{}{
			"esVersion": config.Grafana.ElasticVersion,
			"timeField": "created_time",
		},
	}

	_, err := client.CreateDatasource(options)
	if err != nil {
		log.Println(err)
	}
}

func createDashboard(client *grafana.Client, page Page) {
	f, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Fatal("Can't open dashboard template")
	}
	tpl := string(f)
	t, err := template.New("dashboard").Parse(tpl)
	if err != nil {
		log.Errorln(err)
		return
	}
	data := struct {
		Name string
	}{
		Name: page.Name,
	}

	var options bytes.Buffer
	err = t.Execute(&options, data)
	if err != nil {
		log.Errorln(err)
		return
	}

	_, err = client.CreateDashboard(options.Bytes())
	if err != nil {
		log.Println(err)
	}
}
