Suchbar
=======

![Logo](http://i.imgur.com/I5mjWip.png)

Suchbar edits the public data of any facebook page, making it possible to browse and analyze it.

**Motivation**

For months the PEGIDA-movement (Patriotic Europeans Against the Islamization of the West) has been walking through Dresden. While taciturn in public, PEGIDA's fans do comment very openly on the facebook page. Thus, the idea came up to analyze the contents of the movement and of those persons commenting and liking it.

**Requirements**

- Docker
- Elasticsearch
- Facebook Application (needed for appId, appSecret and accessToken)


**Usage**

	docker run --rm -it 0x46616c6b/suchbar:latest \
	-elastic.host http://elastic:9200 \
	-facebook.app <appID> \
	-facebook.secret <appSecret> \
	-facebook.page <pageID> \
	-facebook.since 24h


**CLI Flags**

	  -elastic.host string
    	the elasticsearch host (default "http://localhost:9200")
      -facebook.app string
    	the app id
      -facebook.limit int
    	the limit for fetching posts per iteration (default 100)
      -facebook.page string
    	the page id
      -facebook.secret string
    	the app secret
      -facebook.since string
    	the earliest date for fetching posts
      -facebook.until string
    	the latest date for fetching posts
