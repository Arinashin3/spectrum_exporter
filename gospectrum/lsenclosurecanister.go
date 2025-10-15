package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type EnclosureCanisterInstance struct {
	Id          string       `json:"id,omitempty"`
	EnclosureId string       `json:"enclosure_id,omitempty"`
	CanisterId  string       `json:"canister_id,omitempty"`
	Status      types.Status `json:"status,omitempty"`
	Type        string       `json:"type,omitempty"`
	NodeId      string       `json:"node_id,omitempty"`
	NodeName    string       `json:"node_name,omitempty"`
}

func (_c *SpectrumClient) GetEnclosureCanister() ([]*EnclosureCanisterInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsEnclosureCanister, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*EnclosureCanisterInstance
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
