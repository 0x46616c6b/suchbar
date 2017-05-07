# Suchbar [![Build Status](https://travis-ci.org/0x46616c6b/suchbar.svg?branch=master)](https://travis-ci.org/0x46616c6b/suchbar)

![Logo](http://i.imgur.com/I5mjWip.png)

**Suchbar edits the public data of any facebook page, making it possible to browse and analyze it.**

## Motivation

For months the PEGIDA-movement (Patriotic Europeans Against the Islamization of the West) has been walking through Dresden. While taciturn in public, PEGIDA's fans do comment very openly on the facebook page. Thus, the idea came up to analyze the contents of the movement and of those persons commenting and liking it.

## Requirements

- Elasticsearch
- Facebook Application
- (Grafana) - optional for dashboards
- (Docker) - optional for runtime

## Usage

        make build
        build/suchbar -facebook.since=48h
        # grafana helper
        build/grafana -config path/to/config -template path/to/dashboard.tpl

### Docker

For more information see the Dockerfile

        docker run -v /path/to/config.yml:/config.yml 0x46616c6b/suchbar

### Suchbar CLI Flags

        -config string
        path to the configuration file
        -facebook.since string
        the earliest date for fetching posts
        -facebook.until string
        the latest date for fetching posts
        - facebook.only string
        fetch only this page id
        
### Grafana CLI Flags

        -config string
        path to the configuration file
        -template string
        path to the template file

## Configuration

        app_id: <facebook app id>
        app_secret: <facebook app secret>
        elastic:
          host: http://localhost:9200
        grafana:
          host: https://localhost:3000
          api_key: <grafana api key>
          elastic_host: http://localhost:9200
          elastic_version: 1
        log_level: info
        pages:
          -
            id: 796885057031701
            name: PEGIDA # used for grafana dashboard
            alias: pegida # alias for elasticsearch (human readable alias)
