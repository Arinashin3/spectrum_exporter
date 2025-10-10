package api

import (
	"bytes"
	"net/http"
)

type SpectrumAPIPath string

const (
	SpectrumAPIPrefix                        = "/rest"
	SpectrumAPIAuth          SpectrumAPIPath = SpectrumAPIPrefix + "/auth"
	SpectrumAPILsEventLog    SpectrumAPIPath = SpectrumAPIPrefix + "/lseventlog"
	SpectrumAPILsFcmap       SpectrumAPIPath = SpectrumAPIPrefix + "/lsfcmap"
	SpectrumAPILsSystem      SpectrumAPIPath = SpectrumAPIPrefix + "/lssystem"
	SpectrumAPILsSystemStats SpectrumAPIPath = SpectrumAPIPrefix + "/lssystemstats"
)

func (_api SpectrumAPIPath) NewRequest(endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", string(_api), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.URL, err = req.URL.Parse(endpoint + string(_api))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}
