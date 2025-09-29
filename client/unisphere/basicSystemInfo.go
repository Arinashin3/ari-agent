package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type BasicSystemInfoInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content BasicSystemInfoContent `json:"content"`
	} `json:"entries"`
}

type BasicSystemInfoContent struct {
	Id                  string `json:"id,omitempty"`
	Model               string `json:"model,omitempty"`
	Name                string `json:"name,omitempty"`
	SoftwareVersion     string `json:"softwareVersion,omitempty"`
	SoftwareFullVersion string `json:"softwareFullVersion,omitempty"`
	ApiVersion          string `json:"apiVersion,omitempty"`
	EarliestApiVersion  string `json:"earliestApiVersion,omitempty"`
}

func (c *UnisphereClient) GetBasicSystemInfo(fields []string) (*BasicSystemInfoInstances, error) {
	path := "/api/types/basicSystemInfo/instances?compact=true"
	if len(fields) != 0 {
		path += "&fields=" + strings.Join(fields, ",")
	}
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data BasicSystemInfoInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
