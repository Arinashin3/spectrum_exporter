package gospectrum

import (
	"encoding/json"
	"errors"
	"net/http"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
	"strings"
)

type EnclosureCanisterSummary struct {
	EnclosureId string       `json:"enclosure_id,omitempty"`
	CanisterId  string       `json:"canister_id,omitempty"`
	Status      types.Status `json:"status,omitempty"`
	Type        string       `json:"type,omitempty"`
	NodeId      string       `json:"node_id,omitempty"`
	NodeName    string       `json:"node_name,omitempty"`
}

type EnclosureCanisterDetail struct {
	EnclosureId         string       `json:"enclosure_id,omitempty"`
	CanisterId          string       `json:"canister_id,omitempty"`
	Status              types.Status `json:"status,omitempty"`
	Type                string       `json:"type,omitempty"`
	NodeId              string       `json:"node_id,omitempty"`
	NodeName            string       `json:"node_name,omitempty"`
	FRUPartNumber       string       `json:"FRU_part_number,omitempty"`
	FRUIdentity         string       `json:"FRU_identity,omitempty"`
	WWNN                string       `json:"WWNN,omitempty"`
	FirmwareLevel       string       `json:"firmware_level,omitempty"`
	Temperature         types.Number `json:"temperature,omitempty"`
	FaultLED            types.Bool   `json:"fault_LED,omitempty"`
	SESStatus           types.Status `json:"SES_status,omitempty"`
	ErrorSequenceNumber types.Number `json:"error_sequence_number,omitempty"`
	SASPort1Status      types.Status `json:"SAS_port1_status,omitempty"`
	FirmwareLevel2      string       `json:"firmware_level_2,omitempty"`
	FirmwareLevel3      string       `json:"firmware_level_3,omitempty"`
	FirmwareLevel4      string       `json:"firmware_level_4,omitempty"`
	FirmwareLevel5      string       `json:"firmware_level_5,omitempty"`
	FirmwareLevel6      string       `json:"firmware_level_6,omitempty"`
}

// GetEnclosureCanisterSummary function is same with **lsenclosurecanister** command.
//
// default is to get summary of all canisters
// when you put **enclosureId**, can get canister's summary of selected enclosure.
//
// if you want to get canister's detail,
// you can use <GetEnclosureCanisterDetail> function.
func (_c *Client) GetEnclosureCanisterSummary(enclosureId string) ([]*EnclosureCanisterSummary, error) {
	// Create Request
	var err error
	var req *http.Request
	var path = api.SpectrumCommandLsEnclosureCanister.String(enclosureId)
	if req, err = _c.newRequest(path, nil); err != nil {
		return nil, err
	}

	// Send & Parsing
	var data []*EnclosureCanisterSummary
	var resp *Response
	if resp, err = _c.send(req); err == nil {
		err = json.Unmarshal(resp.Body, &data)
	}

	return data, err
}

// GetEnclosureCanisterDetail function is same with **lsenclosurecanister -canister <canister_id> <enclosure_id>** command.
// like : lsenclosurecanister -canister 1 1
//
// this function required **[ enclosureId, canisterId ]**
func (_c *Client) GetEnclosureCanisterDetail(enclosureId string, canisterId string) (*EnclosureCanisterDetail, error) {
	// Check Requirement
	var required []string
	if enclosureId == "" {
		required = append(required, "enclosureId")
	}
	if canisterId == "" {
		required = append(required, "canisterId")
	}
	if len(required) != 0 {
		return nil, errors.New("missing required fields: [" + strings.Join(required, ",") + "]")
	}

	// Create Request & Body
	var err error
	var req *http.Request
	var path = api.SpectrumCommandLsEnclosureCanister.String(enclosureId)
	var body = struct {
		Canister string `json:"canister"`
	}{canisterId}
	reqBody, err := json.Marshal(body)
	if req, err = _c.newRequest(path, reqBody); err != nil {
		return nil, err
	}

	// Send Request & Parsing
	var data *EnclosureCanisterDetail
	var resp *Response
	if resp, err = _c.send(req); err == nil {
		err = json.Unmarshal(resp.Body, &data)
	}

	return data, err

}
