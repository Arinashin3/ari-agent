package unisphere

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type MetricInstances struct {
	Base    string    `json:"base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content MetricContent `json:"content,omitempty"`
	}
}

type MetricContent struct {
	Id                    int    `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Path                  string `json:"path,omitempty"`
	Type                  int32  `json:"type,omitempty"`
	Description           string `json:"description,omitempty"`
	IsHistoricalAvailable bool   `json:"isHistoricalAvailable,omitempty"`
	IsRealtimeAvailable   bool   `json:"isRealtimeAvailable,omitempty"`
	UnitDisplayString     string `json:"unitDisplayString,omitempty"`
}

// GetMetric
//
// choose fields : id, name, path, type, description, isHistoricalAvailable, isRealtimeAvailable, unitDisplayString
// choose filterMode : "historical", "realtime", "both", "none"
func (c *Client) GetMetric(fields []string, filterMode string) (*MetricInstances, error) {
	path := "/api/types/metric/instances?compact=true"
	if len(fields) > 0 {
		path += "&fields=" + strings.Join(fields, ",")
	}
	switch filterMode {
	case "historical":
		path += "&filter=isHistoricalAvailable%20eq%20true"
	case "realtime":
		path += "&filter=isRealtimeAvailable%20eq%20true"
	case "both":
		path += "&filter=isHistoricalAvailable%20eq%20true"
		path += "&filter=isRealtimeAvailable%20eq%20true"
	case "none":

	default:
		return nil, errors.New("Unsupported metric filterMode: " + filterMode)

	}

	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data MetricInstances
	err = json.Unmarshal(body, &data)

	return &data, nil
}
