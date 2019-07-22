package grafana

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Dashboard struct {
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	Version int32  `json:"version"`
}

func (c *Client) CreateDashboard(options []byte) (Dashboard, error) {
	dashboard := Dashboard{}

	r := bytes.NewReader(options)
	req, err := c.newRequest("POST", "/api/dashboards/db", r)
	if err != nil {
		return dashboard, err
	}

	res, err := c.Do(req)
	if err != nil {
		return dashboard, err
	}
	if res.StatusCode != 200 {
		return dashboard, errors.New(res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dashboard, err
	}

	err = json.Unmarshal(data, &dashboard)
	if err != nil {
		return dashboard, err
	}
	return dashboard, nil
}
