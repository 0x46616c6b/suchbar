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
const FacebookAPIVersion = "v2.5"

//FacebookLikesLimit defines the Limit for Paging
const FacebookLikesLimit = "500"

//FacebookCommentsLimit defines the Limit for Paging
const FacebookCommentsLimit = FacebookLikesLimit

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
func (ff *FacebookFetcher) GetLikes(postID string) ([]facebook.Result, error) {
	likes, err := ff.fetch(postID, "likes", map[string]string{"limit": FacebookLikesLimit})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for i, like := range likes {
		userID := like["id"].(string)
		likes[i]["post_id"] = fmt.Sprintf("%s", postID)
		likes[i]["user_id"] = fmt.Sprintf("%s", userID)
		likes[i]["id"] = fmt.Sprintf("%s_%s", postID, userID)
	}

	return likes, nil
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

	log.Debugf("/%s/%s?%s", pageID, endpoint, query.Encode())
	result, err := ff.session.Get(fmt.Sprintf("/%s/%s?%s", pageID, endpoint, query.Encode()), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	paging, err := result.Paging(ff.session)
	if err != nil {
		log.Error(err)
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
