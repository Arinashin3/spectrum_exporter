package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
)

type LsFcMapInst struct {
	Id              string                `json:"id,omitempty"`
	Name            string                `json:"name,omitempty"`
	SourceVdiskId   string                `json:"source_vdisk_id,omitempty"`
	SourceVdiskName string                `json:"source_vdisk_name,omitempty"`
	TargetVdiskId   string                `json:"target_vdisk_id,omitempty"`
	TargetVdiskName string                `json:"target_vdisk_name,omitempty"`
	GroupId         string                `json:"group_id,omitempty"`
	GroupName       string                `json:"group_name,omitempty"`
	Status          FlashCopyStatusString `json:"status,omitempty"`
	Progress        string                `json:"progress,omitempty"`
	CopyRate        string                `json:"copy_rate,omitempty"`
	CleanProgress   string                `json:"clean_progress,omitempty"`
	Incremental     string                `json:"incremental,omitempty"`
	PartnerFCId     string                `json:"partner_FC_id,omitempty"`
	PartnerFCName   string                `json:"partner_FC_name,omitempty"`
	Restoring       string                `json:"restoring,omitempty"`
	StartTime       string                `json:"start_time,omitempty"`
	RcControlled    string                `json:"rc_controlled,omitempty"`
}

type FlashCopyStatusString string

const (
	FlashCopyStatusIdleOrCopied FlashCopyStatusString = "idle_or_copied"
	FlashCopyStatusPreparing    FlashCopyStatusString = "preparing"
	FlashCopyStatusPrepared     FlashCopyStatusString = "prepared"
	FlashCopyStatusCopying      FlashCopyStatusString = "copying"
	FlashCopyStatusStopped      FlashCopyStatusString = "stopped"
	FlashCopyStatusStopping     FlashCopyStatusString = "stopping"
	FlashCopyStatusSuspended    FlashCopyStatusString = "suspended"
)

var FlashCopyStatusMap = map[FlashCopyStatusString]int{
	FlashCopyStatusIdleOrCopied: 0,
	FlashCopyStatusPreparing:    1,
	FlashCopyStatusPrepared:     2,
	FlashCopyStatusCopying:      3,
	FlashCopyStatusStopped:      4,
	FlashCopyStatusStopping:     5,
	FlashCopyStatusSuspended:    6,
}

func (_status FlashCopyStatusString) Enum() int {
	return FlashCopyStatusMap[FlashCopyStatusIdleOrCopied]
}

func (_c *SpectrumClient) PostLsFcMap() ([]*LsFcMapInst, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}

	req, err := api.SpectrumAPILsFcmap.NewRequest(_c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*LsFcMapInst
	err = json.Unmarshal(body, &data)

	return data, err
}
