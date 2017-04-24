package storage

import (
	log "github.com/Sirupsen/logrus"
	"github.com/huandu/facebook"
	"gopkg.in/olivere/elastic.v2"
)

const (
	Post    = "post"
	Comment = "comment"
	Like    = "like"
)

//ElasticStorage holds the Elastic Client and the Index
type ElasticStorage struct {
	Client *elastic.Client
}

//NewElasticStorage returns a new ElasticStorage
func NewElasticStorage(host string) *ElasticStorage {
	c, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}

	return &ElasticStorage{
		Client: c,
	}
}

//EnsureAlias creates an alias for an index if their not exists
func (es *ElasticStorage) EnsureAlias(index, alias string) error {
	chk, err := es.Client.Aliases().
		Indices(index, alias).
		Do()
	if err != nil {
		return err
	}

	if len(chk.Indices) == 2 {
		return nil
	}

	_, err = es.Client.Alias().
		Add(index, alias).
		Do()
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

//SaveLikes store all the likes.
func (es *ElasticStorage) SaveLikes(items []facebook.Result, iName string) error {
	return es.save(items, iName, Like)
}

func (es *ElasticStorage) save(items []facebook.Result, iName, tName string) error {
	for _, item := range items {
		_, err := es.Client.Index().
			Index(iName).
			Type(tName).
			Id(item["id"].(string)).
			BodyJson(item).
			Do()

		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
