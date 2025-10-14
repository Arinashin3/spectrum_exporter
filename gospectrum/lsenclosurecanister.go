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
	Type        types.Type   `json:"type,omitempty"`
	NodeId      string       `json:"node_id,omitempty"`
	NodeName    string       `json:"node_name,omitempty"`
}

func (_c *SpectrumClient) GetEnclosureCanister() ([]*EnclosureCanisterInstance, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}

	req, err := api.SpectrumAPILsEnclosureCanister.NewRequest(_c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)
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
