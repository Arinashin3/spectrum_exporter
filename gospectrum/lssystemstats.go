package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"strings"
)

type LsSystemStatsInst struct {
	StatName     string `json:"stat_name"`
	StatCurrent  string `json:"stat_current,omitempty"`
	StatPeak     string `json:"stat_peak,omitempty"`
	StatPeakTime string `json:"stat_peak_time,omitempty"`
}

type lsSystemStatsOptions struct {
	Filtervalue string `json:"filtervalue"`
}

func (_c *SpectrumClient) PostLsSystemStats(filterValues []string) ([]*LsSystemStatsInst, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}

	// Parse Body
	var reqBody []byte
	if len(filterValues) != 0 {
		var opts lsSystemStatsOptions
		opts.Filtervalue = strings.Join(filterValues, ":")
		reqBody, err = json.Marshal(opts)
		if err != nil {
			return nil, err
		}
	}

	// Create Request
	req, err := api.SpectrumAPILsSystemStats.NewRequest(_c.endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)

	// Send
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*LsSystemStatsInst
	err = json.Unmarshal(body, &data)

	return data, err
}
