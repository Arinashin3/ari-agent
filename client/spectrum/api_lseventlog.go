package spectrum

import (
	"encoding/json"
	"time"
)

type LsEventLogInst struct {
	SequenceNumber string `json:"sequence_number,omitempty"`
	LastTimestamp  string `json:"last_timestamp,omitempty"`
	ObjectType     string `json:"object_type,omitempty"`
	ObjectId       string `json:"object_id,omitempty"`
	ObjectName     string `json:"object_name,omitempty"`
	CopyId         string `json:"copy_id,omitempty"`
	Status         string `json:"status,omitempty"`
	Fixed          string `json:"fixed,omitempty"`
	EventId        string `json:"event_id,omitempty"`
	ErrorCode      string `json:"error_code,omitempty"`
	Description    string `json:"description,omitempty"`
}

type LsEventLogRequest struct {
	Filtervalue string `json:"filtervalue,omitempty"`
}

func (c *Client) PostLsEventLog(currentTime time.Time) []*LsEventLogInst {
	ctime := currentTime.Format("060102150405")
	var reqBody LsEventLogRequest
	reqBody.Filtervalue = "last_timestamp>=" + ctime
	jsonReq, _ := json.Marshal(reqBody)
	body, err := c.post("/rest/lseventlog", jsonReq)
	if err != nil {
		c.lastAuth = false
		return nil
	}
	if body == nil {
		return nil
	}
	var data []*LsEventLogInst
	err = json.Unmarshal(body, &data)

	return data
}
