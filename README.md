Suchbar
=======

![Logo](http://i.imgur.com/I5mjWip.png)

Suchbar edits the public data of any facebook page, making it possible to browse and analyze it.

**Motivation**

For months the PEGIDA-movement (Patriotic Europeans Against the Islamization of the West) has been walking through Dresden. While taciturn in public, PEGIDA's fans do comment very openly on the facebook page. Thus, the idea came up to analyze the contents of the movement and of those persons commenting and liking it.

**Requirements**

- (Docker)
- Elasticsearch
- Facebook Application (needed for appId, appSecret and accessToken)

**Usage**

		go build
		./suchbar -facebook.since 48h

**CLI Flags**

		-config string
		path to the configuration file
		-facebook.since string
		the earliest date for fetching posts
		-facebook.until string
		the latest date for fetching posts

**Config**

		app_id: <facebook app id>
		app_secret: <facebook app secret>
		elastic_host: http://localhost:9200
		log_level: info
		pages:
		  -
		    id: 796885057031701
		    alias: pegida
