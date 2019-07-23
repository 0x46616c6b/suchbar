package storage

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/huandu/facebook"
	"gopkg.in/olivere/elastic.v5"
)

const (
	Post       = "post"
	Comment    = "comment"
	Reaction   = "reaction"
	Attachment = "attachment"
)

//ElasticStorage holds the Elastic Client and the Index
type ElasticStorage struct {
	Client *elastic.Client
	Ctx    context.Context
}

//NewElasticStorage returns a new ElasticStorage
func NewElasticStorage(host string) *ElasticStorage {
	c, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}

	return &ElasticStorage{
		Client: c,
		Ctx:    context.Background(),
	}
}

//EnsureAlias creates an alias for an index if their not exists
func (es *ElasticStorage) EnsureAlias(index, alias string) error {
	chk, err := es.Client.Aliases().Do(es.Ctx)
	if err != nil {
		return err
	}

	if len(chk.IndicesByAlias(alias)) == 1 {
		return nil
	}

	_, err = es.Client.Alias().
		Add(index, alias).
		Do(es.Ctx)
	if err != nil {
		return err
	}

	return nil
}

//SavePosts store all the posts
func (es *ElasticStorage) SavePosts(items []facebook.Result, iName string) error {
	return es.save(items, iName, Post)
}

//SaveComments store all the comments
func (es *ElasticStorage) SaveComments(items []facebook.Result, iName string) error {
	return es.save(items, iName, Comment)
}

//SaveReactions store all the reactions.
func (es *ElasticStorage) SaveReactions(items []facebook.Result, iName string) error {
	return es.save(items, iName, Reaction)
}

//SaveAttachments store all the attachments.
func (es *ElasticStorage) SaveAttachments(items []facebook.Result, iName string) error {
	return es.save(items, iName, Attachment)
}

func (es *ElasticStorage) save(items []facebook.Result, iName, tName string) error {
	if len(items) == 0 {
		// nothing to process here
		return nil
	}

	bulkRequest := es.Client.Bulk()

	for _, item := range items {
		req := elastic.NewBulkIndexRequest().
			Index(iName).
			Type(tName).
			Id(item["id"].(string)).
			Doc(item)

		bulkRequest.Add(req)
	}

	start := time.Now()
	bulkResponse, err := bulkRequest.Do(es.Ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"index": iName,
			"type":  tName,
			"items": len(items),
			"error": err,
		}).Error(`Elasticsearch`)

		return err
	}

	log.WithFields(log.Fields{
		"duration": time.Since(start),
		"index":    iName,
		"type":     tName,
		"items":    len(items),
		"created":  len(bulkResponse.Created()),
	}).Debug("Elasticsearch Save")

	return nil
}
