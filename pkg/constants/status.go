package constants

// ItemStatus represents the status of a backlog item
type ItemStatus string

const (
	ItemStatusNew        ItemStatus = "New"
	ItemStatusReady      ItemStatus = "Ready"
	ItemStatusInProgress ItemStatus = "In Progress"
	ItemStatusDone       ItemStatus = "Done"
	ItemStatusArchived   ItemStatus = "Archived"
)

func (s ItemStatus) IsValid() bool {
	switch s {
	case ItemStatusNew, ItemStatusReady, ItemStatusInProgress, ItemStatusDone, ItemStatusArchived:
		return true
	}
	return false
}

// SprintStatus represents the status of a sprint
type SprintStatus string

const (
	SprintStatusPlanning  SprintStatus = "Planning"
	SprintStatusActive    SprintStatus = "Active"
	SprintStatusCompleted SprintStatus = "Completed"
	SprintStatusCancelled SprintStatus = "Cancelled"
)

func (s SprintStatus) IsValid() bool {
	switch s {
	case SprintStatusPlanning, SprintStatusActive, SprintStatusCompleted, SprintStatusCancelled:
		return true
	}
	return false
}
