package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type NodeCanisterInstance struct {
	Id          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	WWNN        string       `json:"WWNN,omitempty"`
	Status      types.Status `json:"status,omitempty"`
	ConfigNode  types.Bool   `json:"config_node,omitempty"`
	Hardware    string       `json:"hardware,omitempty"`
	IscsiName   string       `json:"iscsi_name,omitempty"`
	IscsiAlias  string       `json:"iscsi_alias,omitempty"`
	EnclosureId string       `json:"enclosure_id,omitempty"`
	CanisterId  string       `json:"canister_id,omitempty"`
}

func (_c *SpectrumClient) GetNodeCanister() ([]*NodeCanisterInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsNodeCanister, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*NodeCanisterInstance
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
