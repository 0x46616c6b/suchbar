{
  "dashboard": {
    "annotations": {
      "list": []
    },
    "editable": true,
    "gnetId": null,
    "hideControls": false,
    "links": [],
    "rows": [
      {
            "title": "Dashboard Row",
            "panels": [
              {
                "aliasColors": {},
                "bars": true,
                "datasource": "{{ .Name }}",
                "fill": 1,
                "id": 1,
                "interval": ">1h",
                "legend": {
                  "alignAsTable": true,
                  "avg": true,
                  "current": false,
                  "max": true,
                  "min": true,
                  "show": true,
                  "total": true,
                  "values": true
                },
                "lines": false,
                "linewidth": 1,
                "links": [],
                "nullPointMode": "connected",
                "percentage": false,
                "pointradius": 5,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "span": 5,
                "stack": false,
                "steppedLine": false,
                "targets": [
                  {
                    "bucketAggs": [
                      {
                        "field": "created_time",
                        "id": "2",
                        "settings": {
                          "interval": "auto",
                          "min_doc_count": 0,
                          "trimEdges": 0
                        },
                        "type": "date_histogram"
                      }
                    ],
                    "dsType": "elasticsearch",
                    "metrics": [
                      {
                        "field": "select field",
                        "id": "1",
                        "type": "count"
                      }
                    ],
                    "query": "_type:post",
                    "refId": "A",
                    "timeField": "created_time"
                  }
                ],
                "thresholds": [],
                "timeFrom": null,
                "timeShift": null,
                "title": "Beiträge",
                "tooltip": {
                  "msResolution": false,
                  "shared": true,
                  "sort": 0,
                  "value_type": "individual"
                },
                "type": "graph",
                "xaxis": {
                  "mode": "time",
                  "name": null,
                  "show": true,
                  "values": []
                },
                "yaxes": [
                  {
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                  },
                  {
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                  }
                ]
              },
              {
                "aliasColors": {},
                "bars": true,
                "datasource": "{{ .Name }}",
                "editable": true,
                "error": false,
                "fill": 1,
                "id": 2,
                "interval": ">1h",
                "legend": {
                  "alignAsTable": true,
                  "avg": true,
                  "current": false,
                  "max": true,
                  "min": true,
                  "show": true,
                  "total": true,
                  "values": true
                },
                "lines": false,
                "linewidth": 1,
                "links": [],
                "nullPointMode": "connected",
                "percentage": false,
                "pointradius": 5,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "span": 5,
                "stack": false,
                "steppedLine": false,
                "targets": [
                  {
                    "bucketAggs": [
                      {
                        "field": "created_time",
                        "id": "2",
                        "settings": {
                          "interval": "auto",
                          "min_doc_count": 0,
                          "trimEdges": 0
                        },
                        "type": "date_histogram"
                      }
                    ],
                    "dsType": "elasticsearch",
                    "metrics": [
                      {
                        "field": "select field",
                        "id": "1",
                        "type": "count"
                      }
                    ],
                    "query": "_type:comment",
                    "refId": "A",
                    "timeField": "created_time"
                  }
                ],
                "thresholds": [],
                "timeFrom": null,
                "timeShift": null,
                "title": "Kommentare",
                "tooltip": {
                  "msResolution": false,
                  "shared": true,
                  "sort": 0,
                  "value_type": "individual"
                },
                "type": "graph",
                "xaxis": {
                  "mode": "time",
                  "name": null,
                  "show": true,
                  "values": []
                },
                "yaxes": [
                  {
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                  },
                  {
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                  }
                ]
              },
              {
                "content": "Das Dashboard zeigt Informationen von einer Facebook Seite an. Alle 10 Minuten werden diese Informationen aktualisiert. Die Daten werden über die [Graph API](https://developers.facebook.com/docs/graph-api/reference/page/) von Facebook in einer Suchdatenbank abgespeichert. Der [Quelltext](https://github.com/0x46616c6b/suchbar) des Programms ist Open Source.",
                "editable": true,
                "error": false,
                "id": 6,
                "links": [],
                "mode": "markdown",
                "span": 2,
                "title": "Information",
                "transparent": false,
                "type": "text"
              }
            ],
            "showTitle": false,
            "titleSize": "h6",
            "height": "250px",
            "repeat": null,
            "repeatRowId": null,
            "repeatIteration": null,
            "collapse": false
          },
          {
            "title": "Beiträge",
            "panels": [
              {
                "columns": [
                  {
                    "text": "created_time",
                    "value": "created_time"
                  },
                  {
                    "text": "id",
                    "value": "id"
                  },
                  {
                    "text": "story",
                    "value": "story"
                  },
                  {
                    "text": "message",
                    "value": "message"
                  }
                ],
                "datasource": "{{ .Name }}",
                "editable": true,
                "error": false,
                "fontSize": "100%",
                "height": "400",
                "id": 5,
                "links": [],
                "pageSize": null,
                "scroll": true,
                "showHeader": true,
                "sort": {
                  "col": 0,
                  "desc": true
                },
                "span": 12,
                "styles": [
                  {
                    "dateFormat": "YYYY-MM-DD HH:mm:ss",
                    "pattern": "created_time",
                    "type": "date"
                  },
                  {
                    "colorMode": null,
                    "colors": [
                      "rgba(245, 54, 54, 0.9)",
                      "rgba(237, 129, 40, 0.89)",
                      "rgba(50, 172, 45, 0.97)"
                    ],
                    "decimals": 2,
                    "pattern": "/.*/",
                    "sanitize": true,
                    "thresholds": [],
                    "type": "string",
                    "unit": "short"
                  }
                ],
                "targets": [
                  {
                    "bucketAggs": [],
                    "dsType": "elasticsearch",
                    "metrics": [
                      {
                        "field": "select field",
                        "id": "1",
                        "meta": {},
                        "settings": {},
                        "type": "raw_document"
                      }
                    ],
                    "query": "_type:post",
                    "refId": "A",
                    "timeField": "created_time"
                  }
                ],
                "title": "",
                "transform": "json",
                "transparent": true,
                "type": "table"
              }
            ],
            "showTitle": true,
            "titleSize": "h2",
            "height": "",
            "repeat": null,
            "repeatRowId": null,
            "repeatIteration": null,
            "collapse": false
          },
          {
            "title": "Kommentare",
            "panels": [
              {
                "columns": [
                  {
                    "text": "created_time",
                    "value": "created_time"
                  },
                  {
                    "text": "from.id",
                    "value": "from.id"
                  },
                  {
                    "text": "from.name",
                    "value": "from.name"
                  },
                  {
                    "text": "id",
                    "value": "id"
                  },
                  {
                    "text": "message",
                    "value": "message"
                  }
                ],
                "datasource": "{{ .Name }}",
                "fontSize": "100%",
                "height": "400",
                "id": 3,
                "links": [],
                "pageSize": null,
                "scroll": true,
                "showHeader": true,
                "sort": {
                  "col": 0,
                  "desc": true
                },
                "span": 12,
                "styles": [
                  {
                    "dateFormat": "YYYY-MM-DD HH:mm:ss",
                    "pattern": "created_time",
                    "type": "date"
                  },
                  {
                    "colorMode": null,
                    "colors": [
                      "rgba(245, 54, 54, 0.9)",
                      "rgba(237, 129, 40, 0.89)",
                      "rgba(50, 172, 45, 0.97)"
                    ],
                    "decimals": 2,
                    "pattern": "/.*/",
                    "sanitize": true,
                    "thresholds": [],
                    "type": "string",
                    "unit": "short"
                  }
                ],
                "targets": [
                  {
                    "bucketAggs": [],
                    "dsType": "elasticsearch",
                    "metrics": [
                      {
                        "field": "select field",
                        "id": "1",
                        "meta": {},
                        "settings": {},
                        "type": "raw_document"
                      }
                    ],
                    "query": "_type:comment",
                    "refId": "A",
                    "timeField": "created_time"
                  }
                ],
                "title": "",
                "transform": "json",
                "transparent": true,
                "type": "table"
              }
            ],
            "showTitle": true,
            "titleSize": "h2",
            "height": "",
            "repeat": null,
            "repeatRowId": null,
            "repeatIteration": null,
            "collapse": false
          }
    ],
    "schemaVersion": 13,
    "sharedCrosshair": false,
    "style": "dark",
    "tags": [],
    "templating": {
      "list": []
    },
    "time": {
      "from": "now-30d",
      "to": "now"
    },
    "timepicker": {
      "refresh_intervals": [
        "5s",
        "10s",
        "30s",
        "1m",
        "5m",
        "15m",
        "30m",
        "1h",
        "2h",
        "1d"
      ],
      "time_options": [
        "5m",
        "15m",
        "1h",
        "6h",
        "12h",
        "24h",
        "2d",
        "7d",
        "30d"
      ]
    },
    "timezone": "browser",
    "title": "{{ .Name }}"
  },
  "overwrite": true
}