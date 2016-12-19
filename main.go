package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/0x46616c6b/suchbar/fetcher"
	"github.com/0x46616c6b/suchbar/storage"
	log "github.com/Sirupsen/logrus"
)

var app string
var secret string
var page string
var esHost string
var since string
var until string
var limit int
var logLevel string

func main() {
	fetcher := fetcher.NewFacebookFetcher(app, secret)
	storage := storage.NewElasticStorage(esHost, page)

	items, err := fetcher.GetPosts(page, buildParams())
	if err != nil {
		quit(err)
	}

	err = storage.SavePosts(items)
	if err != nil {
		quit(err)
	}

	log.Debugf("Fetched %d posts", len(items))

	for _, post := range items {
		postID := post["id"].(string)

		log.Debugf("Fetch comments and likes for %s", postID)

		comments, err := fetcher.GetComments(postID)
		if err != nil {
			quit(err)
		}

		err = storage.SaveComments(comments)
		if err != nil {
			quit(err)
		}

		log.Debugf("Fetched %d comments", len(comments))

		likes, err := fetcher.GetLikes(postID)
		if err != nil {
			quit(err)
		}

		err = storage.SaveLikes(likes)
		if err != nil {
			quit(err)
		}

		log.Debugf("Fetched %d likes", len(likes))
	}
}

func init() {

	flag.StringVar(&app, "facebook.app", "", "the app id")
	flag.StringVar(&secret, "facebook.secret", "", "the app secret")
	flag.StringVar(&page, "facebook.page", "", "the page id")
	flag.StringVar(&esHost, "elastic.host", "http://localhost:9200", "the elasticsearch host")
	flag.StringVar(&since, "facebook.since", "", "the earliest date for fetching posts")
	flag.StringVar(&until, "facebook.until", "", "the latest date for fetching posts")
	flag.IntVar(&limit, "facebook.limit", 100, "the limit for fetching posts per iteration")
	flag.StringVar(&logLevel, "log.level", "error", "log level for logrus")
	flag.Parse()

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}

func buildParams() map[string]string {
	p := map[string]string{}

	p["limit"] = strconv.Itoa(limit)

	if since != "" {
		p["since"] = since
	}

	if until != "" {
		p["until"] = until
	}

	return p
}

func quit(err error) {
	log.Fatal(err)
	os.Exit(1)
}
