package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
)

var (
	since  string
	until  string
	limit  int
	config Config
)

func init() {
	flag.StringVar(&since, "facebook.since", "", "the earliest date for fetching posts")
	flag.StringVar(&until, "facebook.until", "", "the latest date for fetching posts")
	flag.IntVar(&limit, "facebook.limit", 100, "the limit for fetching posts per iteration")
	configPath := flag.String("configPath", "config.yml", "path to the configuration file")
	flag.Parse()

	config.read(*configPath)

	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}

func main() {
	r := NewRunner(config)
	r.Run()
}
