package types

// Status
type Status string

const (
	StatusOffline       Status = "offline"
	StatusOnline        Status = "online"
	StatusDegraded      Status = "degraded"
	StatusDegradedPaths Status = "degraded_paths"
	StatusDegradedPorts Status = "degraded_ports"
	StatusExcluded      Status = "excluded"
)

var StatusEnum = map[Status]float64{
	StatusOffline:       0,
	StatusOnline:        1,
	StatusDegraded:      2,
	StatusDegradedPaths: 2.1,
	StatusDegradedPorts: 2.2,
	StatusExcluded:      3,
}

func (_string Status) Enum() float64 {
	return StatusEnum[_string]
}

// RaidStatus
type RaidStatus string

const (
	RaidStatusOffline      RaidStatus = "offline"
	RaidStatusOnline       RaidStatus = "online"
	RaidStatusDegraded     RaidStatus = "degraded"
	RaidStatusSyncing      RaidStatus = "syncing"
	RaidStatusInitializing RaidStatus = "initializing"
	RaidStatusExpanding    RaidStatus = "expanding"
)

var RaidStatusEnum = map[RaidStatus]float64{
	RaidStatusOffline:      0,
	RaidStatusOnline:       1,
	RaidStatusDegraded:     2,
	RaidStatusSyncing:      3,
	RaidStatusInitializing: 4,
	RaidStatusExpanding:    5,
}

func (_status RaidStatus) Enum() float64 {
	return RaidStatusEnum[_status]
}

// FlashCopyStatus
type FlashCopyStatus string

const (
	FlashCopyStatusEnumIdleOrCopied FlashCopyStatus = "idle_or_copied"
	FlashCopyStatusEnumPreparing    FlashCopyStatus = "preparing"
	FlashCopyStatusEnumPrepared     FlashCopyStatus = "prepared"
	FlashCopyStatusEnumCopying      FlashCopyStatus = "copying"
	FlashCopyStatusEnumStopped      FlashCopyStatus = "stopped"
	FlashCopyStatusEnumStopping     FlashCopyStatus = "stopping"
	FlashCopyStatusEnumSuspended    FlashCopyStatus = "suspended"
)

var FlashCopyStatusEnum = map[FlashCopyStatus]float64{
	FlashCopyStatusEnumIdleOrCopied: 0,
	FlashCopyStatusEnumPreparing:    1,
	FlashCopyStatusEnumPrepared:     2,
	FlashCopyStatusEnumCopying:      3,
	FlashCopyStatusEnumStopped:      4,
	FlashCopyStatusEnumStopping:     5,
	FlashCopyStatusEnumSuspended:    6,
}

func (_status FlashCopyStatus) Enum() float64 {
	return FlashCopyStatusEnum[_status]
}
