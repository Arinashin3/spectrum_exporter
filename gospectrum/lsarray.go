package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type ArrayInstance struct {
	MdiskId      string          `json:"mdisk_id,omitempty"`
	MdiskName    string          `json:"mdisk_name,omitempty"`
	Status       types.Status    `json:"status,omitempty"`
	MdiskGrpId   string          `json:"mdiskGrpId,omitempty"`
	MdiskGrpName string          `json:"mdiskGrpName,omitempty"`
	Capacity     types.IEC       `json:"capacity,omitempty"`
	RaidStatus   types.Status    `json:"raidStatus,omitempty"`
	RaidLevel    types.RaidLevel `json:"raid_level"`
	Redundancy   string          `json:"redundancy"`
	StripSize    string          `json:"strip_size,omitempty"`
	Tier         string          `json:"tier,omitempty"`
	Encrypt      types.Bool      `json:"encrypt,omitempty"`
	Distributed  types.Bool      `json:"distributed,omitempty"`
}

func (_c *SpectrumClient) GetArray() ([]*ArrayInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsArray, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*ArrayInstance
	err = json.Unmarshal(body, &data)

	return data, err
}
