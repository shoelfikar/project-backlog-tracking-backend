package constants

// Priority represents the priority level of a backlog item
type Priority string

const (
	PriorityCritical Priority = "Critical"
	PriorityHigh     Priority = "High"
	PriorityMedium   Priority = "Medium"
	PriorityLow      Priority = "Low"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow:
		return true
	}
	return false
}

// ItemType represents the type of a backlog item
type ItemType string

const (
	ItemTypeStory ItemType = "Story"
	ItemTypeBug   ItemType = "Bug"
	ItemTypeTask  ItemType = "Task"
	ItemTypeEpic  ItemType = "Epic"
)

func (t ItemType) IsValid() bool {
	switch t {
	case ItemTypeStory, ItemTypeBug, ItemTypeTask, ItemTypeEpic:
		return true
	}
	return false
}
