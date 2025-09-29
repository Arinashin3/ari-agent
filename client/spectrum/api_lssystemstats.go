package spectrum

import "encoding/json"

type LsSystemStatsInst struct {
	StatName     string `json:"stat_name"`
	StatCurrent  string `json:"stat_current,omitempty"`
	StatPeak     string `json:"stat_peak,omitempty"`
	StatPeakTime string `json:"stat_peak_time,omitempty"`
}

func (c *Client) PostLsSystemStats() []*LsSystemStatsInst {
	body, err := c.post("/rest/lssystemstats", nil)
	if err != nil {
		c.lastAuth = false
		return nil
	}
	if body == nil {
		return nil
	}
	var data []*LsSystemStatsInst
	err = json.Unmarshal(body, &data)

	return data
}
