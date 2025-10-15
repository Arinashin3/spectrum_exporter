package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type VDiskInstance struct {
	Id                  string       `json:"id,omitempty"`
	Name                string       `json:"name,omitempty"`
	IOGroupId           string       `json:"IO_group_id,omitempty"`
	IOGroupName         string       `json:"IO_group_name,omitempty"`
	Status              types.Status `json:"status,omitempty"`
	MdiskGrpId          string       `json:"mdisk_grp_id,omitempty"`
	MdiskGrpName        string       `json:"mdisk_grp_name,omitempty"`
	Capacity            types.IEC    `json:"capacity,omitempty"`
	Type                string       `json:"type,omitempty"`
	FCId                string       `json:"FC_id,omitempty"`
	FCName              string       `json:"FC_name,omitempty"`
	RCId                string       `json:"RC_id,omitempty"`
	RCName              string       `json:"RC_name,omitempty"`
	VdiskUID            string       `json:"vdisk_UID,omitempty"`
	FcMapCount          types.Number `json:"fc_map_count,omitempty"`
	CopyCount           types.Number `json:"copy_count,omitempty"`
	FastWriteState      string       `json:"fast_write_state,omitempty"`
	SeCopyCount         types.Number `json:"se_copy_count,omitempty"`
	RcChange            types.Bool   `json:"rc_change,omitempty"`
	CompressedCopyCount types.Number `json:"compressed_copy_count,omitempty"`
	Formatting          types.Bool   `json:"formatting,omitempty"`
	Encrypt             types.Bool   `json:"encrypt,omitempty"`
	VolumeId            string       `json:"volume_id,omitempty"`
	VolumeName          string       `json:"volume_name,omitempty"`
	Protocol            string       `json:"protocol,omitempty"`
}

func (_c *SpectrumClient) GetVDisk() ([]*VDiskInstance, error) {

	req, err := _c.newRequest(api.SpectrumAPILsVDisk, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*VDiskInstance
	err = json.Unmarshal(body, &data)

	return data, err
}
