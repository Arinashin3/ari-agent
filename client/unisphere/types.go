package unisphere

import (
	"net/http"
	"time"
)

type Client struct {
	endpoint   string
	auth       string
	token      string
	lastAccess bool
	accessTime time.Time
	hc         *http.Client
}

type HealthEnum int

const (
	HealthEnumUnknown        HealthEnum = 0
	HealthEnumOk                        = 5
	HealthEnumOkBut                     = 7
	HealthEnumDegraded                  = 10
	HealthEnumMinor                     = 15
	HealthEnumMajor                     = 20
	HealthEnumCritical                  = 25
	HealthEnumNonRecoverable            = 30
)

type Health struct {
	Value          HealthEnum `json:"value"`
	DescriptionIds []string   `json:"descriptionIds"`
	Descriptions   []string   `json:"descriptions"`
	ResolutionIds  []string   `json:"resolutionIds"`
	Resolutions    []string   `json:"resolutions"`
}
