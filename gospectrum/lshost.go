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

func (_c *Client) GetHost() ([]*HostInstance, error) {
	req, err := _c.newRequest(api.SpectrumCommandLsHost.String(""), nil)
	if err != nil {
		return nil, err
	}
	resp, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*HostInstance
	err = json.Unmarshal(resp.Body, &data)

	return data, err
}
