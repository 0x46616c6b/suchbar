package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/0x46616c6b/suchbar/fetcher"
	"github.com/0x46616c6b/suchbar/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/huandu/facebook"
)

//Runner contains the Storage, Fetcher and Configs
type Runner struct {
	Fetcher *fetcher.FacebookFetcher
	Storage *storage.ElasticStorage
	Config  Config
}

//NewRunner creates an instance of Runner
func NewRunner(c Config) *Runner {
	return &Runner{
		Fetcher: fetcher.NewFacebookFetcher(c.AppID, c.AppSecret),
		Storage: storage.NewElasticStorage(c.ElasticHost),
		Config:  c,
	}
}

//Run starts the Fetcher and stores the Data
func (r *Runner) Run() {
	var wg sync.WaitGroup
	start := time.Now()

	for _, page := range r.Config.Pages {
		wg.Add(1)
		go func(p Page) {
			defer wg.Done()
			r.process(p)
		}(page)
	}

	wg.Wait()
	log.Printf("Fetching took %s", time.Since(start))
}

func (r *Runner) process(p Page) {
	log.Printf("Starting process for page %s", p.Alias)

	items, err := r.Fetcher.GetPosts(p.ID, buildParams())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to fetch posts from "%s"`, p.Alias)
		return
	}
	err = r.Storage.SavePosts(items, p.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to store posts from "%s"`, p.Alias)
	}
	log.Debugf("Fetched %d posts", len(items))

	err = r.Storage.EnsureAlias(p.ID, p.Alias)
	if err != nil {
		log.Error(err)
	}

	r.processPosts(items, p)
}

func (r *Runner) processPosts(items []facebook.Result, p Page) {
	for _, post := range items {
		postID := post["id"].(string)
		log.Debugf("Fetch comments and likes for %s", postID)

		comments, err := r.Fetcher.GetComments(postID)
		if err != nil {
			log.Error(err)
		}

		err = r.Storage.SaveComments(comments, p.ID)
		if err != nil {
			log.Error(err)
		}

		log.Debugf("Fetched %d comments", len(comments))
		likes, err := r.Fetcher.GetLikes(postID)
		if err != nil {
			log.Error(err)
		}

		err = r.Storage.SaveLikes(likes, p.ID)
		if err != nil {
			log.Error(err)
		}

		log.Debugf("Fetched %d likes", len(likes))
	}
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
