package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
	"strings"
)

type SystemStatsInstance struct {
	StatName     string          `json:"stat_name"`
	StatCurrent  string          `json:"stat_current,omitempty"`
	StatPeak     string          `json:"stat_peak,omitempty"`
	StatPeakTime types.Timestamp `json:"stat_peak_time,omitempty"`
}

type SystemStatsOptions struct {
	Filtervalue string `json:"filtervalue"`
}

func (_c *Client) GetSystemStats(filterValues []string) ([]*SystemStatsInstance, error) {

	// Parse Body
	var reqBody []byte
	var err error
	if len(filterValues) != 0 {
		var opts SystemStatsOptions
		opts.Filtervalue = strings.Join(filterValues, ":")
		reqBody, err = json.Marshal(opts)
		if err != nil {
			return nil, err
		}
	}

	// Create Request
	req, err := _c.newRequest(api.SpectrumCommandLsSystemStats.String(""), reqBody)
	if err != nil {
		return nil, err
	}

	// Send
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*SystemStatsInstance
	err = json.Unmarshal(body.Body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
