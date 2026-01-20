package service

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/dto/response"
	"sprint-backlog/internal/models"
	"sprint-backlog/internal/repository"
	"sprint-backlog/pkg/constants"
)

var (
	ErrBacklogItemNotFound = errors.New("backlog item not found")
	ErrInvalidItemType     = errors.New("invalid item type")
	ErrInvalidPriority     = errors.New("invalid priority")
	ErrInvalidStatus       = errors.New("invalid status")
)

type BacklogService interface {
	Create(req *request.CreateBacklogItemRequest, userID uuid.UUID) (*response.BacklogItemResponse, error)
	GetByID(id uuid.UUID) (*response.BacklogItemResponse, error)
	GetAll(params *request.BacklogQueryParams) (*response.BacklogListResponse, error)
	Update(id uuid.UUID, req *request.UpdateBacklogItemRequest, userID uuid.UUID) (*response.BacklogItemResponse, error)
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status constants.ItemStatus, userID uuid.UUID) (*response.BacklogItemResponse, error)
	UpdatePriority(id uuid.UUID, priority constants.Priority, userID uuid.UUID) (*response.BacklogItemResponse, error)
	AddLabel(id uuid.UUID, label string, userID uuid.UUID) (*response.BacklogItemResponse, error)
	RemoveLabel(id uuid.UUID, label string, userID uuid.UUID) (*response.BacklogItemResponse, error)
	AddComment(id uuid.UUID, content string, userID uuid.UUID) (*response.ItemHistoryResponse, error)
	GetHistory(id uuid.UUID) ([]response.ItemHistoryResponse, error)
}

type backlogService struct {
	backlogRepo repository.BacklogRepository
	historyRepo repository.ItemHistoryRepository
}

func NewBacklogService(backlogRepo repository.BacklogRepository, historyRepo repository.ItemHistoryRepository) BacklogService {
	return &backlogService{
		backlogRepo: backlogRepo,
		historyRepo: historyRepo,
	}
}

func (s *backlogService) Create(req *request.CreateBacklogItemRequest, userID uuid.UUID) (*response.BacklogItemResponse, error) {
	// Validate type
	if !req.Type.IsValid() {
		return nil, ErrInvalidItemType
	}

	// Validate priority
	if !req.Priority.IsValid() {
		return nil, ErrInvalidPriority
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = constants.ItemStatusNew
	}
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	// Get max position
	maxPos, err := s.backlogRepo.GetMaxPosition(req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Create item
	item := &models.BacklogItem{
		ProjectID:   req.ProjectID,
		SprintID:    req.SprintID,
		CreatedByID: userID,
		Title:       strings.TrimSpace(req.Title),
		Type:        req.Type,
		Priority:    req.Priority,
		Status:      status,
		StoryPoints: req.StoryPoints,
		Labels:      req.Labels,
		Position:    maxPos + 1,
	}

	// Set description if provided
	if req.Description != "" {
		desc := strings.TrimSpace(req.Description)
		item.Description = &desc
	}

	if err := s.backlogRepo.Create(item); err != nil {
		return nil, err
	}

	// Record history
	s.recordHistory(item.ID, userID, constants.ItemActionCreated, nil, nil, nil, nil)

	// Fetch created item with relations
	created, err := s.backlogRepo.GetByID(item.ID)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogItemResponse(created), nil
}

func (s *backlogService) GetByID(id uuid.UUID) (*response.BacklogItemResponse, error) {
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	return response.ToBacklogItemResponse(item), nil
}

func (s *backlogService) GetAll(params *request.BacklogQueryParams) (*response.BacklogListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Build filters
	filters := repository.BacklogFilters{
		Search: params.Search,
		Page:   params.Page,
		Limit:  params.Limit,
	}

	// Parse type filter
	for _, t := range params.Type {
		itemType := constants.ItemType(t)
		if itemType.IsValid() {
			filters.Type = append(filters.Type, itemType)
		}
	}

	// Parse priority filter
	for _, p := range params.Priority {
		priority := constants.Priority(p)
		if priority.IsValid() {
			filters.Priority = append(filters.Priority, priority)
		}
	}

	// Parse status filter
	for _, s := range params.Status {
		status := constants.ItemStatus(s)
		if status.IsValid() {
			filters.Status = append(filters.Status, status)
		}
	}

	// Parse sprint filter
	if params.SprintID != "" {
		if params.SprintID == "none" {
			nilID := uuid.Nil
			filters.SprintID = &nilID
		} else if sprintID, err := uuid.Parse(params.SprintID); err == nil {
			filters.SprintID = &sprintID
		}
	}

	// Parse labels filter
	filters.Labels = params.Labels

	items, total, err := s.backlogRepo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogListResponse(items, total, params.Page, params.Limit), nil
}

func (s *backlogService) Update(id uuid.UUID, req *request.UpdateBacklogItemRequest, userID uuid.UUID) (*response.BacklogItemResponse, error) {
	// Get existing item
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	// Track changes for history
	changes := make(map[string][2]interface{})

	// Update fields if provided
	if req.Title != "" && req.Title != item.Title {
		changes["title"] = [2]interface{}{item.Title, req.Title}
		item.Title = strings.TrimSpace(req.Title)
	}

	if req.Description != "" {
		oldDesc := ""
		if item.Description != nil {
			oldDesc = *item.Description
		}
		if req.Description != oldDesc {
			changes["description"] = [2]interface{}{oldDesc, req.Description}
			desc := strings.TrimSpace(req.Description)
			item.Description = &desc
		}
	}

	if req.Type != "" && req.Type != item.Type {
		if !req.Type.IsValid() {
			return nil, ErrInvalidItemType
		}
		changes["type"] = [2]interface{}{item.Type, req.Type}
		item.Type = req.Type
	}

	if req.Priority != "" && req.Priority != item.Priority {
		if !req.Priority.IsValid() {
			return nil, ErrInvalidPriority
		}
		changes["priority"] = [2]interface{}{item.Priority, req.Priority}
		item.Priority = req.Priority
	}

	if req.Status != "" && req.Status != item.Status {
		if !req.Status.IsValid() {
			return nil, ErrInvalidStatus
		}
		changes["status"] = [2]interface{}{item.Status, req.Status}
		item.Status = req.Status
	}

	if req.StoryPoints != nil {
		changes["story_points"] = [2]interface{}{item.StoryPoints, req.StoryPoints}
		item.StoryPoints = req.StoryPoints
	}

	if req.Labels != nil {
		changes["labels"] = [2]interface{}{item.Labels, req.Labels}
		item.Labels = req.Labels
	}

	if req.SprintID != nil {
		changes["sprint_id"] = [2]interface{}{item.SprintID, req.SprintID}
		item.SprintID = req.SprintID
	}

	if err := s.backlogRepo.Update(item); err != nil {
		return nil, err
	}

	// Record history for each change
	for field, vals := range changes {
		oldVal, _ := json.Marshal(vals[0])
		newVal, _ := json.Marshal(vals[1])
		s.recordHistory(id, userID, constants.ItemActionUpdated, &field, datatypes.JSON(oldVal), datatypes.JSON(newVal), nil)
	}

	// Fetch updated item with relations
	updated, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogItemResponse(updated), nil
}

func (s *backlogService) Delete(id uuid.UUID) error {
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrBacklogItemNotFound
	}

	return s.backlogRepo.Delete(id)
}

func (s *backlogService) UpdateStatus(id uuid.UUID, status constants.ItemStatus, userID uuid.UUID) (*response.BacklogItemResponse, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	// Get current item
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	oldStatus := item.Status

	// Update status
	if err := s.backlogRepo.UpdateStatus(id, status); err != nil {
		return nil, err
	}

	// Record history
	oldVal, _ := json.Marshal(oldStatus)
	newVal, _ := json.Marshal(status)
	field := "status"
	s.recordHistory(id, userID, constants.ItemActionStatusChanged, &field, datatypes.JSON(oldVal), datatypes.JSON(newVal), nil)

	// Fetch updated item
	updated, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogItemResponse(updated), nil
}

func (s *backlogService) UpdatePriority(id uuid.UUID, priority constants.Priority, userID uuid.UUID) (*response.BacklogItemResponse, error) {
	if !priority.IsValid() {
		return nil, ErrInvalidPriority
	}

	// Get current item
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	oldPriority := item.Priority

	// Update priority
	if err := s.backlogRepo.UpdatePriority(id, priority); err != nil {
		return nil, err
	}

	// Record history
	oldVal, _ := json.Marshal(oldPriority)
	newVal, _ := json.Marshal(priority)
	field := "priority"
	s.recordHistory(id, userID, constants.ItemActionPriorityChanged, &field, datatypes.JSON(oldVal), datatypes.JSON(newVal), nil)

	// Fetch updated item
	updated, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogItemResponse(updated), nil
}

func (s *backlogService) AddLabel(id uuid.UUID, label string, userID uuid.UUID) (*response.BacklogItemResponse, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, errors.New("label cannot be empty")
	}

	// Check if item exists
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	// Add label
	if err := s.backlogRepo.AddLabel(id, label); err != nil {
		return nil, err
	}

	// Record history
	newVal, _ := json.Marshal(label)
	s.recordHistory(id, userID, constants.ItemActionLabelAdded, nil, nil, datatypes.JSON(newVal), nil)

	// Fetch updated item
	updated, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogItemResponse(updated), nil
}

func (s *backlogService) RemoveLabel(id uuid.UUID, label string, userID uuid.UUID) (*response.BacklogItemResponse, error) {
	// Check if item exists
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	// Remove label
	if err := s.backlogRepo.RemoveLabel(id, label); err != nil {
		return nil, err
	}

	// Record history
	oldVal, _ := json.Marshal(label)
	s.recordHistory(id, userID, constants.ItemActionLabelRemoved, nil, datatypes.JSON(oldVal), nil, nil)

	// Fetch updated item
	updated, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToBacklogItemResponse(updated), nil
}

func (s *backlogService) AddComment(id uuid.UUID, content string, userID uuid.UUID) (*response.ItemHistoryResponse, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("comment cannot be empty")
	}

	// Check if item exists
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	// Create history entry for comment
	history := &models.ItemHistory{
		ItemID:  id,
		UserID:  userID,
		Action:  constants.ItemActionCommentAdded,
		Comment: &content,
	}

	if err := s.historyRepo.Create(history); err != nil {
		return nil, err
	}

	return response.ToItemHistoryResponse(history), nil
}

func (s *backlogService) GetHistory(id uuid.UUID) ([]response.ItemHistoryResponse, error) {
	// Check if item exists
	item, err := s.backlogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	histories, err := s.historyRepo.GetByItemID(id)
	if err != nil {
		return nil, err
	}

	return response.ToItemHistoryListResponse(histories), nil
}

// recordHistory is a helper function to record item history
func (s *backlogService) recordHistory(itemID, userID uuid.UUID, action constants.ItemAction, field *string, oldValue, newValue datatypes.JSON, comment *string) {
	history := &models.ItemHistory{
		ItemID:       itemID,
		UserID:       userID,
		Action:       action,
		FieldChanged: field,
		OldValue:     oldValue,
		NewValue:     newValue,
		Comment:      comment,
	}
	s.historyRepo.Create(history)
}
