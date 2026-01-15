package constants

// ItemAction represents actions that can be performed on a backlog item
type ItemAction string

const (
	ItemActionCreated            ItemAction = "Created"
	ItemActionUpdated            ItemAction = "Updated"
	ItemActionStatusChanged      ItemAction = "StatusChanged"
	ItemActionPriorityChanged    ItemAction = "PriorityChanged"
	ItemActionSprintAssigned     ItemAction = "SprintAssigned"
	ItemActionSprintRemoved      ItemAction = "SprintRemoved"
	ItemActionCommentAdded       ItemAction = "CommentAdded"
	ItemActionLabelAdded         ItemAction = "LabelAdded"
	ItemActionLabelRemoved       ItemAction = "LabelRemoved"
	ItemActionDescriptionUpdated ItemAction = "DescriptionUpdated"
)

// SprintAction represents actions that can be performed on a sprint
type SprintAction string

const (
	SprintActionCreated     SprintAction = "Created"
	SprintActionStarted     SprintAction = "Started"
	SprintActionItemAdded   SprintAction = "ItemAdded"
	SprintActionItemRemoved SprintAction = "ItemRemoved"
	SprintActionItemMoved   SprintAction = "ItemMoved"
	SprintActionCompleted   SprintAction = "Completed"
	SprintActionCancelled   SprintAction = "Cancelled"
)
