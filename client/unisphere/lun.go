package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type LunInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content LunContent `json:"content"`
	} `json:"entries"`
}

type LunContent struct {
	Id                     string      `json:"id,omitempty"`
	Health                 Health      `json:"health,omitempty"`
	Name                   string      `json:"name,omitempty"`
	Description            string      `json:"description,omitempty"`
	Type                   LunTypeEnum `json:"type,omitempty"`
	SizeTotal              int64       `json:"sizeTotal,omitempty"`
	SizeUsed               int64       `json:"sizeUsed,omitempty"`
	SizeAllocated          int64       `json:"sizeAllocated,omitempty"`
	SizePreallocated       int64       `json:"sizePreallocated,omitempty"`
	SizeAllocatedTotal     int64       `json:"sizeAllocatedTotal,omitempty"`
	DataReductionSizeSaved int64       `json:"dataReductionSizeSaved,omitempty"`
	DataReductionPercent   int64       `json:"dataReductionPercent,omitempty"`
	DataReductionRatio     int64       `json:"dataReductionRatio,omitempty"`
	IsThinEnabled          bool        `json:"isThinEnabled,omitempty"`
	IsDataReductionEnabled bool        `json:"isDataReductionEnabled,omitempty"`
	IsAdvancedDedupEnabled bool        `json:"isAdvancedDedupEnabled,omitempty"`
	Wwn                    string      `json:"wwn,omitempty"`
}

type LunTypeEnum int

const (
	LunTypeGenericStorageStandalone LunTypeEnum = iota + 1
	LunTypeStandalone
	LunTypeVmWareISCSI
)

func (c *UnisphereClient) GetLun(fields []string) (*LunInstances, error) {
	path := "/api/types/lun/instances?compact=true"
	if len(fields) != 0 {
		path += "&fields=" + strings.Join(fields, ",")
	}
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var data LunInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
