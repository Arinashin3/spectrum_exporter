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

func (_c *Client) GetMDisk() ([]*MDiskInstance, error) {
	req, err := _c.newRequest(api.SpectrumCommandLsMDisk.String(""), nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*MDiskInstance
	err = json.Unmarshal(body.Body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
