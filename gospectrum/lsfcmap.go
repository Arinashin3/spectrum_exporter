package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type FlashCopyMapInstance struct {
	Id              string                `json:"id,omitempty"`
	Name            string                `json:"name,omitempty"`
	SourceVdiskId   string                `json:"source_vdisk_id,omitempty"`
	SourceVdiskName string                `json:"source_vdisk_name,omitempty"`
	TargetVdiskId   string                `json:"target_vdisk_id,omitempty"`
	TargetVdiskName string                `json:"target_vdisk_name,omitempty"`
	GroupId         string                `json:"group_id,omitempty"`
	GroupName       string                `json:"group_name,omitempty"`
	Status          types.FlashCopyStatus `json:"status,omitempty"`
	Progress        types.Number          `json:"progress,omitempty"`
	CopyRate        types.Number          `json:"copy_rate,omitempty"`
	CleanProgress   types.Number          `json:"clean_progress,omitempty"`
	Incremental     types.Bool            `json:"incremental,omitempty"`
	PartnerFCId     string                `json:"partner_FC_id,omitempty"`
	PartnerFCName   string                `json:"partner_FC_name,omitempty"`
	Restoring       types.Bool            `json:"restoring,omitempty"`
	StartTime       types.Timestamp       `json:"start_time,omitempty"`
	RcControlled    types.Bool            `json:"rc_controlled,omitempty"`
}

func (_c *SpectrumClient) GetFlashCopyMap() ([]*FlashCopyMapInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsFcMap, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, nil
	}
	var data []*FlashCopyMapInstance
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
