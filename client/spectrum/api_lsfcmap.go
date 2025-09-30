package spectrum

import (
	"encoding/json"
	"strconv"
	"time"
)

type LsFcMapInst struct {
	Id              string          `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	SourceVdiskId   string          `json:"source_vdisk_id,omitempty"`
	SourceVdiskName string          `json:"source_vdisk_name,omitempty"`
	TargetVdiskId   string          `json:"target_vdisk_id,omitempty"`
	TargetVdiskName string          `json:"target_vdisk_name,omitempty"`
	GroupId         string          `json:"group_id,omitempty"`
	GroupName       string          `json:"group_name,omitempty"`
	Status          FlashCopyStatus `json:"status,omitempty"`
	Progress        String          `json:"progress,omitempty"`
	CopyRate        String          `json:"copy_rate,omitempty"`
	CleanProgress   String          `json:"clean_progress,omitempty"`
	Incremental     string          `json:"incremental,omitempty"`
	PartnerFCId     string          `json:"partner_FC_id,omitempty"`
	PartnerFCName   string          `json:"partner_FC_name,omitempty"`
	Restoring       string          `json:"restoring,omitempty"`
	StartTime       String          `json:"start_time,omitempty"`
	RcControlled    string          `json:"rc_controlled,omitempty"`
}

func (c *Client) PostLsFcMap() []*LsFcMapInst {
	body, err := c.post("/rest/lsfcmap", nil)
	if err != nil {
		c.lastAuth = false
		return nil
	}
	if body == nil {
		return nil
	}
	var data []*LsFcMapInst
	err = json.Unmarshal(body, &data)

	return data
}

type FlashCopyStatus string

const (
	FlashCopyStatusEnumIdleOrCopied FlashCopyStatus = "idle_or_copied"
	FlashCopyStatusEnumPreparing    FlashCopyStatus = "preparing"
	FlashCopyStatusEnumPrepared     FlashCopyStatus = "prepared"
	FlashCopyStatusEnumCopying      FlashCopyStatus = "copying"
	FlashCopyStatusEnumStopped      FlashCopyStatus = "stopped"
	FlashCopyStatusEnumStopping     FlashCopyStatus = "stopping"
	FlashCopyStatusEnumSuspended    FlashCopyStatus = "suspended"
)

func (_fc FlashCopyStatus) Float64() float64 {
	switch _fc {
	case FlashCopyStatusEnumIdleOrCopied:
		return 0.0
	case FlashCopyStatusEnumCopying:
		return 1.0
	case FlashCopyStatusEnumPreparing:
		return 2.0
	case FlashCopyStatusEnumPrepared:
		return 3.0
	case FlashCopyStatusEnumStopped:
		return 4.0
	case FlashCopyStatusEnumStopping:
		return 5.0
	case FlashCopyStatusEnumSuspended:
		return 6.0
	default:
		return -1.0
	}
}

type String string

func (_s String) Float64() float64 {
	f, err := strconv.ParseFloat(string(_s), 64)
	if err != nil {
		f = -999
	}
	return f
}

func (_s String) Time2Float64() float64 {
	t, _ := time.Parse("060102150405", string(_s))
	return float64(t.Unix())
}
