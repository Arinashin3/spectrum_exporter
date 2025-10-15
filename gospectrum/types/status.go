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
	StatusService       Status = "service"
	StatusFlushing      Status = "flushing"
	StatusPending       Status = "pending"
	StatusAdding        Status = "adding"
	StatusDeleting      Status = "deleting"
	StatusSpare         Status = "spare"
	StatusOnlineSpare   Status = "online_spare"
	StatusSyncing       Status = "syncing"
	StatusInitializing  Status = "initializing"
	StatusExpanding     Status = "expanding"
)

var StatusEnum = map[Status]float64{
	StatusOffline:       0,
	StatusOnline:        1,
	StatusDegraded:      2,
	StatusDegradedPaths: 3,
	StatusDegradedPorts: 4,
	StatusExcluded:      5,
	StatusService:       6,
	StatusFlushing:      7,
	StatusPending:       8,
	StatusAdding:        9,
	StatusDeleting:      10,
	StatusSpare:         11,
	StatusOnlineSpare:   12,
	StatusSyncing:       13,
	StatusInitializing:  14,
	StatusExpanding:     15,
}

func (_string Status) Enum() float64 {
	return StatusEnum[_string]
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
