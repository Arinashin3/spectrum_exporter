package gospectrum

import (
	"encoding/json"
	"net/http"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type EnclosureSummary struct {
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

// GetEnclosureSummary function is same with 'lsenclosure' command.
//
// if you want to get enclosure details, put enclosure_id.
// default is to get summary information of all enclosure.
func (_c *Client) GetEnclosureSummary(enclosureId string) ([]*EnclosureSummary, error) {
	// Create Request
	var err error
	var req *http.Request
	var path = api.SpectrumCommandLsEnclosure.String(enclosureId)
	if req, err = _c.newRequest(path, nil); err != nil {
		return nil, err
	}

	// Send & Parsing
	var data []*EnclosureSummary
	var resp *Response
	if resp, err = _c.send(req); err == nil {
		err = json.Unmarshal(resp.Body, &data)
	}

	return data, err
}
