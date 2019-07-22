package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Datasource struct {
	ID                int    `json:"id"`
	OrgID             int    `json:"orgId"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	TypeLogoURL       string `json:"typeLogoUrl"`
	Access            string `json:"access"`
	URL               string `json:"url"`
	Password          string `json:"password"`
	User              string `json:"user"`
	Database          string `json:"database"`
	BasicAuth         bool   `json:"basicAuth"`
	BasicAuthUser     string `json:"basicAuthUser"`
	BasicAuthPassword string `json:"basicAuthPassword"`
	WithCredentials   bool   `json:"withCredentials"`
	IsDefault         bool   `json:"isDefault"`
}

func (c *Client) CreateDatasource(options map[string]interface{}) (Datasource, error) {
	datasource := Datasource{}
	body, err := json.Marshal(options)
	if err != nil {
		return datasource, err
	}

	r := bytes.NewReader(body)
	req, err := c.newRequest("POST", "/api/datasources", r)
	if err != nil {
		return datasource, err
	}

	res, err := c.Do(req)
	if err != nil {
		return datasource, err
	}
	if res.StatusCode != 200 {
		if res.Status == "409 Conflict" {
			return datasource, errors.New("Datasource already exists")
		}
		return datasource, errors.New(res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return datasource, err
	}

	err = json.Unmarshal(data, &datasource)
	if err != nil {
		return datasource, err
	}
	return datasource, nil
}
