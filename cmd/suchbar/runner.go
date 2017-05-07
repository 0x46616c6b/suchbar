package main

import (
	"runtime"
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
		Storage: storage.NewElasticStorage(c.Elastic.Host),
		Config:  c,
	}
}

//Run starts the Fetcher and stores the Data
func (r *Runner) Run() {
	var wg sync.WaitGroup
	start := time.Now()

	for _, page := range r.Config.Pages {
		// skip pages when only argument set and not equal to actual page
		if only != "" && page.ID != only {
			return
		}

		wg.Add(1)
		go func(p Page) {
			defer wg.Done()
			r.process(p)
		}(page)
	}

	wg.Wait()
	log.Printf("Fetching data from all pages took %s", time.Since(start))
}

func (r *Runner) process(p Page) {
	log.Printf("Starting process for page %s", p.Alias)
	start := time.Now()

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
	log.Printf(`Fetching data from "%s" took %s`, p.ID, time.Since(start))
}

type work struct {
	item facebook.Result
}

func (r *Runner) processPosts(items []facebook.Result, p Page) {
	c := make(chan work, 10)
	var wg sync.WaitGroup

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range c {
				postID := w.item["id"].(string)
				r.processComments(postID, p)
				r.processLikes(postID, p)
			}
			return
		}()
	}

	for _, item := range items {
		c <- work{item}
	}

	close(c)
	wg.Wait()
}

func (r *Runner) processComments(postID string, p Page) {
	comments, err := r.Fetcher.GetComments(postID)
	if err != nil {
		log.Error(err)
	}
	err = r.Storage.SaveComments(comments, p.ID)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Fetched %d comments", len(comments))
}

func (r *Runner) processLikes(postID string, p Page) {
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
