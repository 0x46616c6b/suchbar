# Suchbar [![Build Status](https://travis-ci.com/0x46616c6b/suchbar.svg?branch=master)](https://travis-ci.com/0x46616c6b/suchbar)

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

## Datastructure

### Elastic

**Posts**

        {
            "_index": "796885057031701",
            "_type": "post",
            "_id": "796885057031701_1381651581888376",
            "_source": {
               "created_time": "2017-05-25T19:04:59+0000",
               "id": "796885057031701_1381651581888376",
               "message": "#PEGIDA \n\n+++ Flüchtlinge mit Waffen und Drogen festgenommen +++\n\nMittlerweile haben wir quasi eine stehende Armee im Land und das ist nicht die Bundeswehr. . . . \n\n+++ Die Berliner Polizei hat am Mittwoch bei mehreren Razzien in Neukölln, Köpenick, Mariendorf und Zehlendorf neun Flüchtlinge aus dem Irak und Syrien festgenommen. Die Ermittler fanden Betäubungsmittel und Waffen. +++"
            }
        }

**Comments**

        {
            "_index": "796885057031701",
            "_type": "comment",
            "_id": "1381214465265421_1381222908597910",
            "_source": {
               "created_time": "2017-05-25T11:32:46+0000",
               "from": {
                  "id": "136020320105***",
                  "name": "J***** L***"
               },
               "id": "1381214465265421_1381222908597910",
               "message": "\"Service\" - \"La Brigade franco-allemande (BFA)\" http://www.france-allemagne.fr/La-Brigade-franco-allemande-BFA.html  (!Oben links auf die deutsche Fahne klicken! Dann ist es in deutsch!) ;-)"
            }
        }

**Reactions**

        {
            "_index": "796885057031701",
            "_type": "reaction",
            "_id": "796885057031701_1381647395222128_***",
            "_source": {
               "id": "796885057031701_1381647395222128_***",
               "name": "R***** U***",
               "post_id": "796885057031701_1381647395222128",
               "type": "LIKE",
               "user_id": "151594811851***"
            }
        }

**Attachments**

        {
            "_index": "796885057031701",
            "_type": "attachment",
            "_id": "796885057031701_1381472058572995",
            "_score": 1,
            "_source": {
               "description": "Eine 22-jährige Studentin ist am Dienstag von einem mutmaßlichen Sexualtäter angegangen worden. Die junge Frau war kurz nach 20 Uhr auf einem",
               "id": "796885057031701_1381472058572995",
               "media": {
                  "image": {
                     "height": 400,
                     "src": "https://external.xx.fbcdn.net/safe_image.php?d=AQAaGFyNEYjZmadb&w=720&h=720&url=http%3A%2F%2Fwww.stadtzeitung.de%2Fresources%2Fmediadb%2F2017%2F05%2F24%2F85332_web.jpg&cfs=1&_nc_hash=AQCbZn4VdvykaRON",
                     "width": 400
                  }
               },
               "target": {
                  "url": "https://l.facebook.com/l.php?u=http%3A%2F%2Fwww.stadtzeitung.de%2Faugsburg-nordost%2Fblaulicht%2Fjoggerin-wird-an-berliner-allee-ueberfallen-und-sexuell-angegangen-d26812.html%3Fcp%3DKurationsbox&h=ATPLMQ-PPn04D12k54Mow7DxhDM-TRfYi49FDIwgp91Eff-sKd0LjVU2BEss7aBAXBL0CJa9H0zSadGlTF6MukejnKAHhll_iA&s=1&enc=AZOyEUPn6DRlVHC1NFN0ZWQZ5qEkCNu6iZhRb-1FxuFi2JWQRzz9PQIbICjFE-mV0uC-Eao71pbUkpsroQ61SChn"
               },
               "title": "Joggerin wird an Berliner Allee überfallen und sexuell angegangen",
               "type": "share",
               "url": "https://l.facebook.com/l.php?u=http%3A%2F%2Fwww.stadtzeitung.de%2Faugsburg-nordost%2Fblaulicht%2Fjoggerin-wird-an-berliner-allee-ueberfallen-und-sexuell-angegangen-d26812.html%3Fcp%3DKurationsbox&h=ATPLMQ-PPn04D12k54Mow7DxhDM-TRfYi49FDIwgp91Eff-sKd0LjVU2BEss7aBAXBL0CJa9H0zSadGlTF6MukejnKAHhll_iA&s=1&enc=AZOyEUPn6DRlVHC1NFN0ZWQZ5qEkCNu6iZhRb-1FxuFi2JWQRzz9PQIbICjFE-mV0uC-Eao71pbUkpsroQ61SChn"
            }
        }
