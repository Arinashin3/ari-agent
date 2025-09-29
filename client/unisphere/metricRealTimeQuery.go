package unisphere

import (
	"encoding/json"
	"strconv"
	"time"
)

type MetricRealTimeQueryInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Content struct {
		Id int64 `json:"id"`
	} `json:"content"`
}

func (c *UnisphereClient) PostMetricRealTimeQuery(paths []string, interval time.Duration) (string, error) {
	path := "/api/types/metricRealTimeQuery/instances"
	var reqBodySt struct {
		Paths    []string `json:"paths"`
		Interval int      `json:"interval"`
	}
	reqBodySt.Paths = paths
	reqBodySt.Interval = int(interval.Seconds())

	reqBody, err := json.Marshal(reqBodySt)
	if err != nil {
		return "", err
	}

	body, err := c.post(path, reqBody)
	if err != nil {
		return "", err
	}

	var data MetricRealTimeQueryInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(data.Content.Id)), nil
}
