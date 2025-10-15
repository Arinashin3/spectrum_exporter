package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type DriveInstance struct {
	Id                      string          `json:"id,omitempty"`
	Status                  types.Status    `json:"status,omitempty"`
	ErrorSequenceNumber     string          `json:"error_sequence_number,omitempty"`
	Use                     string          `json:"use,omitempty"`
	UID                     string          `json:"UID,omitempty"`
	TechType                string          `json:"tech_type,omitempty"`
	ReplacementDate         types.Timestamp `json:"replacement_date,omitempty"`
	Capacity                types.IEC       `json:"capacity,omitempty"`
	BlockSize               types.Number    `json:"block_size,omitempty"`
	VendorId                string          `json:"vendor_id,omitempty"`
	ProductId               string          `json:"product_id,omitempty"`
	FRUPartNumber           string          `json:"FRU_part_number,omitempty"`
	FRUIdentity             string          `json:"FRU_identity,omitempty"`
	RPM                     types.Number    `json:"RPM,omitempty"`
	FirmwareLevel           string          `json:"Firmware_level,omitempty"`
	FPGALevel               string          `json:"FPGA_level,omitempty"`
	MdiskId                 string          `json:"mdisk_id,omitempty"`
	MdiskName               string          `json:"mdisk_name,omitempty"`
	MemberId                string          `json:"member_id,omitempty"`
	EnclosureId             string          `json:"enclosure_id,omitempty"`
	SlotId                  string          `json:"slot_id,omitempty"`
	NodeName                string          `json:"node_name,omitempty"`
	NodeId                  string          `json:"node_id,omitempty"`
	QuorumId                string          `json:"quorum_id,omitempty"`
	Port1Status             string          `json:"port_1_status,omitempty"`
	Port2Status             string          `json:"port_2_status,omitempty"`
	AutoManage              string          `json:"auto_manage,omitempty"`
	DriveClassId            string          `json:"drive_class_id,omitempty"`
	WriteEnduranceUsed      string          `json:"write_endurance_used,omitempty"`
	WriteEnduranceUsageRate string          `json:"write_endurance_usage_rate,omitempty"`
	TransportProtocol       string          `json:"transport_protocol,omitempty"`
	Compressed              types.Bool      `json:"compressed,omitempty"`
	PhysicalCapacity        types.IEC       `json:"physical_capacity,omitempty"`
	PhysicalUsedCapacity    types.IEC       `json:"physical_used_capacity,omitempty"`
	EffectiveUsedCapacity   types.IEC       `json:"effective_used_capacity,omitempty"`
	DateOfManufacture       types.Timestamp `json:"date_of_manufacture,omitempty"`
}

func (_c *SpectrumClient) GetDrive() ([]*DriveInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsDrive, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*DriveInstance
	err = json.Unmarshal(body, &data)

	return data, err
}
