package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type MgmtInterfaceInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content MgmtInterfaceContent `json:"content"`
	} `json:"entries"`
}

type MgmtInterfaceContent struct {
	Id              string `json:"id,omitempty"`
	ConfigMode      int    `json:"configMode,omitempty"`
	IpAddress       string `json:"ipAddress,omitempty"`
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	Netmask         string `json:"netmask,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
}

func (c *UnisphereClient) GetMgmtInterface(fields []string) (*MgmtInterfaceInstances, error) {
	path := "/api/types/mgmtInterface/instances?compact=true"
	if len(fields) != 0 {
		path += "&fields=" + strings.Join(fields, ",")
	}
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data MgmtInterfaceInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
