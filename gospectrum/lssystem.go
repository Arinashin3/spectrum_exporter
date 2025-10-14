package gospectrum

import (
	"encoding/json"
	"spectrum_exporter/gospectrum/api"
	"spectrum_exporter/gospectrum/types"
)

type SystemInstance struct {
	Id                              string       `json:"id,omitempty"`
	Name                            string       `json:"name,omitempty"`
	Location                        string       `json:"location,omitempty"`
	TotalMdiskCapacity              types.IEC    `json:"total_mdisk_capacity,omitempty"`
	SpaceInMdiskGrps                types.IEC    `json:"space_in_mdisk_grps,omitempty"`
	SpaceAllocatedToVdisks          types.IEC    `json:"space_allocated_to_vdisks,omitempty"`
	TotalFreeSpace                  types.IEC    `json:"total_free_space,omitempty"`
	TotalVdiskcopyCapacity          types.IEC    `json:"total_vdiskcopy_capacity,omitempty"`
	TotalUsedCapacity               types.IEC    `json:"total_used_capacity,omitempty"`
	TotalOverallocation             string       `json:"total_overallocation,omitempty"`
	TotalVdiskCapacity              types.IEC    `json:"total_vdisk_capacity,omitempty"`
	TotalAllocatedExtentCapacity    types.IEC    `json:"total_allocated_extent_capacity,omitempty"`
	StatisticsStatus                types.Bool   `json:"statistics_status,omitempty"`
	StatisticsFrequency             types.Number `json:"statistics_frequency,omitempty"`
	ClusterLocale                   string       `json:"cluster_locale,omitempty"`
	TimeZone                        string       `json:"time_zone,omitempty"`
	CodeLevel                       string       `json:"code_level,omitempty"`
	ConsoleIP                       string       `json:"console_IP,omitempty"`
	IdAlias                         string       `json:"id_alias,omitempty"`
	RelationshipBandwidthLimit      types.Number `json:"relationship_bandwidth_limit,omitempty"`
	Layer                           string       `json:"layer,omitempty"`
	RcBufferSize                    types.Number `json:"rc_buffer_size,omitempty"`
	CompressionActive               types.Bool   `json:"compression_active,omitempty"`
	CompressionVirtualCapacity      types.IEC    `json:"compression_virtual_capacity,omitempty"`
	CompressionCompressedCapacity   types.IEC    `json:"compression_compressed_capacity,omitempty"`
	CompressionUncompressedCapacity types.IEC    `json:"compression_uncompressed_capacity,omitempty"`
	TotalDriveRawCapacity           types.IEC    `json:"total_drive_raw_capacity,omitempty"`
	VdiskProtectionTime             types.Number `json:"vdisk_protection_time,omitempty"`
	VdiskProtectionEnabled          types.Bool   `json:"vdisk_protection_enabled,omitempty"`
	ProductName                     string       `json:"product_name,omitempty"`
	MaxReplicationDelay             types.Number `json:"max_replication_delay,omitempty"`
	TotalReclaimableCapacity        types.IEC    `json:"total_reclaimable_capacity,omitempty"`
	PhysicalCapacity                types.IEC    `json:"physical_capacity,omitempty"`
	PhysicalFreeCapacity            types.IEC    `json:"physical_free_capacity,omitempty"`
	UsedCapacityBeforeReduction     types.IEC    `json:"used_capacity_before_reduction,omitempty"`
	UsedCapacityAfterReduction      types.IEC    `json:"used_capacity_after_reduction,omitempty"`
	OverheadCapacity                types.IEC    `json:"overhead_capacity,omitempty"`
	DeduplicationCapacitySaving     types.IEC    `json:"deduplication_capacity_saving,omitempty"`
}

func (_c *SpectrumClient) GetSystem() (*SystemInstance, error) {
	// Try Login
	err := _c.login()
	if err != nil {
		return nil, err
	}

	var data *SystemInstance
	req, err := api.SpectrumAPILsSystem.NewRequest(_c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", _c.token)
	body, err := _c.send(req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &data)

	return data, err
}
