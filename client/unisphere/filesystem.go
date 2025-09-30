package unisphere

import (
	"encoding/json"
	"strings"
	"time"
)

type FilesystemInstances struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Entries []struct {
		Content FilesystemContent `json:"content"`
	} `json:"entries"`
}

type FilesystemContent struct {
	Id            string             `json:"id,omitempty"`
	Health        Health             `json:"health,omitempty"`
	Name          string             `json:"name,omitempty"`
	Description   string             `json:"description,omitempty"`
	Type          FilesystemTypeEnum `json:"type,omitempty"`
	SizeTotal     int64              `json:"sizeTotal,omitempty"`
	SizeUsed      int64              `json:"sizeUsed,omitempty"`
	SizeAllocated int64              `json:"sizeAllocated,omitempty"`
}

type FilesystemTypeEnum int

const (
	FilesystemTypeFilesystem FilesystemTypeEnum = iota + 1
	FilesystemTypeVMware
)

func (c *UnisphereClient) GetFilesystem(fields []string) (*FilesystemInstances, error) {
	joinedFields := strings.Join(fields, ",")
	body, err := c.get("/api/types/filesystem/instances?compact=true&fields=" + joinedFields)
	if err != nil {
		return nil, err
	}

	var data FilesystemInstances
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
