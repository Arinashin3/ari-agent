package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type SystemCapacityInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content SystemCapacityContent `json:"content"`
	} `json:"entries"`
}

type SystemCapacityContent struct {
	Id                     string  `json:"id,omitempty"`
	SizeFree               int64   `json:"sizeFree,omitempty"`
	SizeTotal              int64   `json:"sizeTotal,omitempty"`
	SizeUsed               int64   `json:"sizeUsed,omitempty"`
	SizePreallocated       int64   `json:"sizePreallocated,omitempty"`
	DataReductionSizeSaved int64   `json:"dataReductionSizeSaved,omitempty"`
	DataReductionPercent   int64   `json:"dataReductionPercent,omitempty"`
	DataReductionRatio     float64 `json:"dataReductionRatio,omitempty"`
	SizeSubscribed         int64   `json:"sizeSubscribed,omitempty"`
	TotalLogicalSize       int64   `json:"totalLogicalSize,omitempty"`
	ThinSavingRatio        float64 `json:"thinSavingRatio,omitempty"`
	SnapsSavingsRatio      float64 `json:"snapsSavingsRatio,omitempty"`
	OverallEfficiencyRatio float64 `json:"overallEfficiencyRatio,omitempty"`
	Tiers                  []struct {
		TierType  int64 `json:"tierType,omitempty"`
		SizeFree  int64 `json:"sizeFree,omitempty"`
		SizeTotal int64 `json:"sizeTotal,omitempty"`
		SizeUsed  int64 `json:"sizeUsed,omitempty"`
	} `json:"titers,omitempty"`
}

func (c *UnisphereClient) GetSystemCapacity(fields []string) (*SystemCapacityInstances, error) {
	joinedFields := strings.Join(fields, ",")
	body, err := c.get("/api/types/systemCapacity/instances?compact=true&fields=" + joinedFields)
	if err != nil {
		return nil, err
	}

	var data SystemCapacityInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
