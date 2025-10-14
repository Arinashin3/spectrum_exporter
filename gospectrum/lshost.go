package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type HostInstance struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	PortCount  types.Number `json:"port_count"`
	Type       string       `json:"type"`
	IogrpCount types.Number `json:"iogrp_count"`
	Status     types.Status `json:"status"`
	Protocol   string       `json:"protocol"`
}

func (_c *SpectrumClient) GetHost() ([]*HostInstance, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}

	req, err := api.SpectrumAPILsHost.NewRequest(_c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*HostInstance
	err = json.Unmarshal(body, &data)

	return data, err
}
