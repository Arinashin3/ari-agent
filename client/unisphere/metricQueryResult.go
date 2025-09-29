package unisphere

import (
	"encoding/json"
	"errors"
	"time"
)

type MetricQueryResultInstance struct {
	Base    string    `json:"base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content struct {
			QueryId int         `json:"queryId,omitempty"`
			Path    string      `json:"path,omitempty"`
			Values  interface{} `json:"values,omitempty"`
		} `json:"content,omitempty"`
	} `json:"entries"`
}

func (c *UnisphereClient) GetMetricQueryResult(queryId string) (*MetricQueryResultInstance, error) {
	path := "/api/types/metricQueryResult/instances?compact=true"
	if queryId == "" {
		return nil, errors.New("queryId is required")
	}
	path += "&filter=queryId%20eq%20" + queryId
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data MetricQueryResultInstance
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
