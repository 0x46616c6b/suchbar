package storage

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/huandu/facebook"
	"gopkg.in/olivere/elastic.v2"
)

//ElasticStorage holds the Elastic Client and the Index
type ElasticStorage struct {
	Client *elastic.Client
	Index  string
}

//NewElasticStorage returns a new ElasticStorage
func NewElasticStorage(host, index string) *ElasticStorage {
	c, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &ElasticStorage{
		Client: c,
		Index:  index,
	}
}

//SavePosts store all the posts
func (es *ElasticStorage) SavePosts(items []facebook.Result) error {
	return es.save(items, "post")
}

//SaveComments store all the comments
func (es *ElasticStorage) SaveComments(items []facebook.Result) error {
	return es.save(items, "comment")
}

//SaveLikes store all the likes.
func (es *ElasticStorage) SaveLikes(items []facebook.Result) error {
	return es.save(items, "like")
}

func (es *ElasticStorage) save(items []facebook.Result, typeName string) error {
	for _, item := range items {
		_, err := es.Client.Index().
			Index(es.Index).
			Type(typeName).
			Id(item["id"].(string)).
			BodyJson(item).
			Do()

		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
