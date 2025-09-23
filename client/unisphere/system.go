package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type SystemInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content SystemContent `json:"content"`
	} `json:"entries"`
}

type SystemContent struct {
	Id     string `json:"id,omitempty"`
	Health struct {
		Value         int      `json:"value,omitempty"`
		DesciptionIds []string `json:"descriptionIds,omitempty"`
		Descriptions  []string `json:"descriptions,omitempty"`
		ResolutionIds []string `json:"resolutionIds,omitempty"`
		Resolutions   []string `json:"resolutions,omitempty"`
	} `json:"health,omitempty"`
	Name                         string `json:"name,omitempty"`
	Model                        string `json:"model,omitempty"`
	SerialNumber                 string `json:"serialNumber,omitempty"`
	UuidBase                     int32  `json:"uuidBase,omitempty"`
	InternalModel                string `json:"internalModel,omitempty"`
	Platform                     string `json:"platform,omitempty"`
	IsAllFlash                   bool   `json:"isAllFlash,omitempty"`
	MacAddress                   string `json:"macAddress,omitempty"`
	IsEULAAccount                bool   `json:"isEULAAccount,omitempty"`
	IsUpgradeComplete            bool   `json:"isUpgradeComplete,omitempty"`
	isAutoFailbackEnabled        bool   `json:"isAutoFailbackEnabled,omitempty"`
	CurrentPower                 int32  `json:"currentPower,omitempty"`
	AvgPower                     int32  `json:"avgPower,omitempty"`
	SupportedUpgradeModels       []int  `json:"supportedUpgradeModels,omitempty"`
	IsRemoteSysInterfaceAutoPair bool   `json:"isRemoteSysInterfaceAutoPair,omitempty"`
}

func (c *Client) GetSystem(fields []string) (*SystemInstances, error) {
	path := "/api/types/system/instances?compact=true"
	if len(fields) != 0 {
		path += "&fields=" + strings.Join(fields, ",")
	}

	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data SystemInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
