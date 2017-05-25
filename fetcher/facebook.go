package fetcher

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/huandu/facebook"
)

//FacebookAPIVersion holds the Facebook API Version
const FacebookAPIVersion = "v2.6"

//FacebookCommentsLimit defines the Limit for Paging
const FacebookCommentsLimit = "500"

//FacebookReactionsLimit defines the Limit for Paging
const FacebookReactionsLimit = FacebookCommentsLimit

//FacebookFetcher holds the session from facebook
type FacebookFetcher struct {
	session *facebook.Session
}

//NewFacebookFetcher return a new FacebookFetcher
func NewFacebookFetcher(appID, appSecret string) *FacebookFetcher {
	app := facebook.New(appID, appSecret)
	facebook.Version = FacebookAPIVersion

	return &FacebookFetcher{session: app.Session(app.AppAccessToken())}
}

//GetPosts return all the posts for a pageID
func (ff *FacebookFetcher) GetPosts(pageID string, params map[string]string) ([]facebook.Result, error) {
	return ff.fetch(pageID, "posts", params)
}

//GetComments returns all the comments for a postID
func (ff *FacebookFetcher) GetComments(postID string) ([]facebook.Result, error) {
	return ff.fetch(postID, "comments", map[string]string{"limit": FacebookCommentsLimit})
}

//GetLikes returns all the likes for a postID
func (ff *FacebookFetcher) GetReactions(postID string) ([]facebook.Result, error) {
	reactions, err := ff.fetch(postID, "reactions", map[string]string{"limit": FacebookReactionsLimit})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for i, reaction := range reactions {
		userID := reaction["id"].(string)
		reactions[i]["post_id"] = fmt.Sprintf("%s", postID)
		reactions[i]["user_id"] = fmt.Sprintf("%s", userID)
		reactions[i]["id"] = fmt.Sprintf("%s_%s", postID, userID)
	}

	return reactions, nil
}

//GetAttachments returns all the attachments for a postID
func (ff *FacebookFetcher) GetAttachments(postID string) ([]facebook.Result, error) {
	attachments, err := ff.fetch(postID, "attachments", map[string]string{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for i := range attachments {
		attachments[i]["id"] = postID
	}

	return attachments, nil
}

func (ff *FacebookFetcher) fetch(pageID string, endpoint string, params map[string]string) ([]facebook.Result, error) {
	query := url.Values{}

	if val, ok := params["limit"]; ok {
		query.Add("limit", val)
	}

	if val, ok := params["since"]; ok {
		query.Add("since", strconv.FormatInt(calcTime(val).Unix(), 10))
	}
	if val, ok := params["until"]; ok {
		query.Add("until", strconv.FormatInt(calcTime(val).Unix(), 10))
	}

	result, err := ff.session.Get(fmt.Sprintf("/%s/%s?%s", pageID, endpoint, query.Encode()), nil)
	if err != nil {
		return nil, err
	}

	paging, err := result.Paging(ff.session)
	if err != nil {
		return nil, err
	}

	items := paging.Data()

	for noMore, err := paging.Next(); !noMore; noMore, err = paging.Next() {
		if err != nil {
			log.Error(err)
			break
		}

		items = append(items, paging.Data()...)
	}

	return items, nil
}

func calcTime(s string) time.Time {
	now := time.Now()
	dur, _ := time.ParseDuration(s)

	return now.Add(-dur)
}
