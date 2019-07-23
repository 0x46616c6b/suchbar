package main

import (
	"runtime"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/huandu/facebook"

	"github.com/0x46616c6b/suchbar/fetcher"
	"github.com/0x46616c6b/suchbar/storage"
)

const POST = "post"
const COMMENT = "comment"
const REACTION = "reaction"
const ATTACHMENT = "attachment"

//Runner contains the Storage, Fetcher and Configs
type Runner struct {
	Fetcher *fetcher.FacebookFetcher
	Storage *storage.ElasticStorage
	Config  Config
	Results map[string]PageResult
}

type PageResult struct {
	Duration    time.Duration
	Posts       int
	Comments    int
	Reactions   int
	Attachments int
}

//NewRunner creates an instance of Runner
func NewRunner(c Config) *Runner {
	// initialize PageResult map
	results := map[string]PageResult{}
	for _, page := range c.Pages {
		results[page.ID] = PageResult{}
	}

	return &Runner{
		Fetcher: fetcher.NewFacebookFetcher(c.AppID, c.AppSecret),
		Storage: storage.NewElasticStorage(c.Elastic.Host),
		Config:  c,
		Results: results,
	}
}

//Run starts the Fetcher and stores the Data
func (r *Runner) Run() {
	log.WithFields(log.Fields{
		"since": since,
		"until": until,
		"limit": limit,
	}).Info("Start fetching pages")
	start := time.Now()
	wg := sync.WaitGroup{}

	if log.GetLevel() == log.DebugLevel {
		yellow := color.New(color.FgYellow).SprintFunc()
		go func() {
			for {
				log.WithFields(log.Fields{
					"NumGoroutine": runtime.NumGoroutine(),
					"NumCgoCall":   runtime.NumCgoCall(),
				}).Debugf(`%s`, yellow(`Runtime Stats`))
				time.Sleep(time.Second * 2)
			}
		}()
	}

	for _, page := range r.Config.Pages {
		// skip pages when only argument set and not equal to actual page
		if only != "" && page.ID != only {
			continue
		}

		wg.Add(1)
		go func(p Page) {
			defer wg.Done()
			r.process(p)
		}(page)
	}

	wg.Wait()
	log.Infof("Duration (Total): %s", time.Since(start))
}

func (r *Runner) process(p Page) {
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
		}).Errorf(`Failed to save posts from "%s"`, p.Alias)
		return
	}
	r.setPostsCounter(p, len(items))

	err = r.Storage.EnsureAlias(p.ID, p.Alias)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to create alias for "%s"`, p.Alias)
	}

	r.processPosts(items, p)
	r.setDuration(p, time.Since(start))

	r.printResults(p)
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
				r.processAttachments(postID, p)
				r.processComments(postID, p)
				r.processReactions(postID, p)
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
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to fetch comments from "%s"`, p.Alias)
		return
	}
	err = r.Storage.SaveComments(comments, p.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to save comments from "%s"`, p.Alias)
		return
	}
	r.setCommentsCounter(p, len(comments))
}

func (r *Runner) processReactions(postID string, p Page) {
	reactions, err := r.Fetcher.GetReactions(postID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to fetch reactions from "%s"`, p.Alias)
		return
	}
	err = r.Storage.SaveReactions(reactions, p.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to save reactions from "%s"`, p.Alias)
		return
	}
	r.setReactionsCounter(p, len(reactions))
}

func (r *Runner) processAttachments(postID string, p Page) {
	attachments, err := r.Fetcher.GetAttachments(postID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to fetch attachments from "%s"`, p.Alias)
		return
	}
	err = r.Storage.SaveAttachments(attachments, p.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  p.ID,
		}).Errorf(`Failed to save attachments from "%s"`, p.Alias)
		return
	}
	r.setAttachmentsCounter(p, len(attachments))
}

func (r *Runner) setPostsCounter(p Page, c int) {
	r.setResultsCounter(p, POST, c)
}

func (r *Runner) setCommentsCounter(p Page, c int) {
	r.setResultsCounter(p, COMMENT, c)
}

func (r *Runner) setReactionsCounter(p Page, c int) {
	r.setResultsCounter(p, REACTION, c)
}

func (r *Runner) setAttachmentsCounter(p Page, c int) {
	r.setResultsCounter(p, ATTACHMENT, c)
}

func (r *Runner) setDuration(p Page, d time.Duration) {
	if _, ok := r.Results[p.ID]; ok {
		results := r.Results[p.ID]
		results.Duration = d

		r.Results[p.ID] = results
	}
}

func (r *Runner) printResults(p Page) {
	if _, ok := r.Results[p.ID]; ok {
		results := r.Results[p.ID]
		fields := log.Fields{
			"duration":    results.Duration,
			"posts":       results.Posts,
			"comments":    results.Comments,
			"reactions":   results.Reactions,
			"attachments": results.Attachments,
		}

		red := color.New(color.FgRed).SprintFunc()
		log.WithFields(fields).Infof("Statistics for %s", red(p.Alias))
	}
}

func (r *Runner) setResultsCounter(p Page, name string, count int) {
	if _, ok := r.Results[p.ID]; ok {
		results := r.Results[p.ID]
		switch name {
		case POST:
			results.Posts = results.Posts + count
			break
		case COMMENT:
			results.Comments = results.Comments + count
			break
		case REACTION:
			results.Reactions = results.Reactions + count
			break
		case ATTACHMENT:
			results.Attachments = results.Attachments + count
		}

		r.Results[p.ID] = results
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
