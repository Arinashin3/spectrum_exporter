package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type MDiskInstance struct {
	Id              string       `json:"id,omitempty"`
	Name            string       `json:"name,omitempty"`
	Status          types.Status `json:"status,omitempty"`
	Mode            types.Mode   `json:"mode,omitempty"`
	MdiskGrpId      string       `json:"mdiskGrpId,omitempty"`
	MdiskGrpName    string       `json:"mdiskGrpName,omitempty"`
	Capacity        types.IEC    `json:"capacity,omitempty"`
	ControllerName  string       `json:"controller_name"`
	UID             string       `json:"UID,omitempty"`
	Distributed     types.Bool   `json:"distributed,omitempty"`
	Dedupe          types.Bool   `json:"dedupe,omitempty"`
	OverProvisioned types.Bool   `json:"over_provisioned,omitempty"`
	SupportsUnmap   types.Bool   `json:"supports_unmap,omitempty"`
}

func (_c *SpectrumClient) GetMDisk() ([]*MDiskInstance, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}

	req, err := api.SpectrumAPILsMdisk.NewRequest(_c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*MDiskInstance
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
