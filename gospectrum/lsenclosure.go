package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type EnclosureInstance struct {
	Id               string       `json:"id,omitempty"`
	Status           types.Status `json:"status,omitempty"`
	Type             string       `json:"type,omitempty"`
	Managed          types.Bool   `json:"managed,omitempty"`
	IOGroupId        string       `json:"IO_group_id,omitempty"`
	IOGroupName      string       `json:"IO_group_name,omitempty"`
	ProductMTM       string       `json:"product_MTM,omitempty"`
	SerialNumber     string       `json:"serial_number,omitempty"`
	TotalCanisters   types.Number `json:"total_canisters,omitempty"`
	OnlineCanisters  types.Number `json:"online_canisters,omitempty"`
	TotalPSUs        types.Number `json:"total_PSUs,omitempty"`
	OnlinePSUs       types.Number `json:"online_PSUs,omitempty"`
	DriveSlots       types.Number `json:"drive_slots,omitempty"`
	TotalFanModules  types.Number `json:"total_fan_modules,omitempty"`
	OnlineFanModules types.Number `json:"online_fan_modules,omitempty"`
	TotalSems        types.Number `json:"total_sems,omitempty"`
	OnlineSems       types.Number `json:"online_sems,omitempty"`
}

func (_c *SpectrumClient) GetEnclosure() ([]*EnclosureInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsEnclosure, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*EnclosureInstance
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
