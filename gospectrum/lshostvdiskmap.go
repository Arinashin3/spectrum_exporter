package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
)

type HostVdiskMapInstance struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	SCSIId          string `json:"SCSI_id,omitempty"`
	HostClusterId   string `json:"host_cluster_id,omitempty"`
	HostClusterName string `json:"host_cluster_name,omitempty"`
	VdiskId         string `json:"vdisk_id,omitempty"`
	VdiskName       string `json:"vdisk_name,omitempty"`
	VdiskUID        string `json:"vdisk_UID,omitempty"`
	IOGroupId       string `json:"IO_group_id,omitempty"`
	IOGroupName     string `json:"IO_group_name,omitempty"`
	MappingType     string `json:"mapping_type,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
}

func (_c *SpectrumClient) GetHostVdiskMap() ([]*HostVdiskMapInstance, error) {
	req, err := _c.newRequest(api.SpectrumAPILsHostVDiskMap, nil)
	if err != nil {
		return nil, err
	}
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	var data []*HostVdiskMapInstance
	err = json.Unmarshal(body, &data)

	return data, err
}
