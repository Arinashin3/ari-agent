package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type EventInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content EventContent `json:"content,omitempty"`
	} `json:"entries"`
}

type EventContent struct {
	Id           string    `json:"id,omitempty"`
	Severity     int64     `json:"severity,omitempty"`
	CreationTime time.Time `json:"creationTime,omitempty"`
	MessageId    string    `json:"messageId,omitempty"`
	Message      string    `json:"message,omitempty"`
	Source       string    `json:"source,omitempty"`
	Arguments    []string  `json:"arguments,omitempty"`
}

func (c *UnisphereClient) GetEvent(fields []string, createTime time.Time) (*EventInstances, error) {
	path := "/api/types/event/instances?compact=true"
	if len(fields) != 0 {
		path += "&fields=" + strings.Join(fields, ",")
	}

	if !createTime.IsZero() {
		path += "&filter=creationTime%20gt%20\"" + createTime.UTC().Format("2006-01-02T15:04:05Z") + "\""
	}

	var body []byte
	var err error
	body, err = c.get(path)
	if err != nil {
		return nil, err
	}

	var data EventInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
