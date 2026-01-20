package service

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/dto/response"
	"sprint-backlog/internal/models"
	"sprint-backlog/internal/repository"
	"sprint-backlog/pkg/constants"
)

var (
	ErrSprintNotFound       = errors.New("sprint not found")
	ErrSprintAlreadyActive  = errors.New("there is already an active sprint in this project")
	ErrSprintNotPlanning    = errors.New("sprint must be in planning status to start")
	ErrSprintNotActive      = errors.New("sprint must be in active status to complete or cancel")
	ErrInvalidSprintStatus  = errors.New("invalid sprint status")
	ErrInvalidDateRange     = errors.New("end date must be after start date")
	ErrItemAlreadyInSprint  = errors.New("item is already in this sprint")
	ErrItemNotInSprint      = errors.New("item is not in this sprint")
)

type SprintService interface {
	Create(req *request.CreateSprintRequest, userID uuid.UUID) (*response.SprintResponse, error)
	GetByID(id uuid.UUID) (*response.SprintResponse, error)
	GetAll(params *request.SprintQueryParams) (*response.SprintListResponse, error)
	GetWithItems(id uuid.UUID) (*response.SprintWithItemsResponse, error)
	GetActive(projectID uuid.UUID) (*response.SprintResponse, error)
	Update(id uuid.UUID, req *request.UpdateSprintRequest, userID uuid.UUID) (*response.SprintResponse, error)
	Delete(id uuid.UUID) error
	Start(id uuid.UUID, userID uuid.UUID) (*response.SprintResponse, error)
	Complete(id uuid.UUID, userID uuid.UUID) (*response.SprintResponse, error)
	Cancel(id uuid.UUID, userID uuid.UUID) (*response.SprintResponse, error)
	AddItem(sprintID uuid.UUID, itemID uuid.UUID, userID uuid.UUID) (*response.SprintWithItemsResponse, error)
	RemoveItem(sprintID uuid.UUID, itemID uuid.UUID, userID uuid.UUID) (*response.SprintWithItemsResponse, error)
	GetHistory(id uuid.UUID) ([]response.SprintHistoryResponse, error)
	GetReport(id uuid.UUID) (*response.SprintReportResponse, error)
}

type sprintService struct {
	sprintRepo        repository.SprintRepository
	sprintHistoryRepo repository.SprintHistoryRepository
	backlogRepo       repository.BacklogRepository
	itemHistoryRepo   repository.ItemHistoryRepository
}

func NewSprintService(
	sprintRepo repository.SprintRepository,
	sprintHistoryRepo repository.SprintHistoryRepository,
	backlogRepo repository.BacklogRepository,
	itemHistoryRepo repository.ItemHistoryRepository,
) SprintService {
	return &sprintService{
		sprintRepo:        sprintRepo,
		sprintHistoryRepo: sprintHistoryRepo,
		backlogRepo:       backlogRepo,
		itemHistoryRepo:   itemHistoryRepo,
	}
}

func (s *sprintService) Create(req *request.CreateSprintRequest, userID uuid.UUID) (*response.SprintResponse, error) {
	// Validate date range
	if !req.EndDate.After(req.StartDate) {
		return nil, ErrInvalidDateRange
	}

	// Create sprint
	sprint := &models.Sprint{
		ProjectID:   req.ProjectID,
		CreatedByID: userID,
		Name:        strings.TrimSpace(req.Name),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      constants.SprintStatusPlanning,
	}

	// Set goal if provided
	if req.Goal != "" {
		goal := strings.TrimSpace(req.Goal)
		sprint.Goal = &goal
	}

	if err := s.sprintRepo.Create(sprint); err != nil {
		return nil, err
	}

	// Record history
	s.recordSprintHistory(sprint.ID, userID, nil, constants.SprintActionCreated, nil, nil)

	// Fetch created sprint with relations
	created, err := s.sprintRepo.GetByID(sprint.ID)
	if err != nil {
		return nil, err
	}

	return response.ToSprintResponse(created), nil
}

func (s *sprintService) GetByID(id uuid.UUID) (*response.SprintResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	return response.ToSprintResponse(sprint), nil
}

func (s *sprintService) GetAll(params *request.SprintQueryParams) (*response.SprintListResponse, error) {
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
	filters := repository.SprintFilters{
		Page:  params.Page,
		Limit: params.Limit,
	}

	// Parse project filter
	if params.ProjectID != "" {
		if projectID, err := uuid.Parse(params.ProjectID); err == nil {
			filters.ProjectID = &projectID
		}
	}

	// Parse status filter
	for _, st := range params.Status {
		status := constants.SprintStatus(st)
		if status.IsValid() {
			filters.Status = append(filters.Status, status)
		}
	}

	sprints, total, err := s.sprintRepo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	return response.ToSprintListResponse(sprints, total, params.Page, params.Limit), nil
}

func (s *sprintService) GetWithItems(id uuid.UUID) (*response.SprintWithItemsResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	items, err := s.sprintRepo.GetItemsBySprintID(id)
	if err != nil {
		return nil, err
	}

	return response.ToSprintWithItemsResponse(sprint, items), nil
}

func (s *sprintService) GetActive(projectID uuid.UUID) (*response.SprintResponse, error) {
	sprint, err := s.sprintRepo.GetActive(projectID)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	return response.ToSprintResponse(sprint), nil
}

func (s *sprintService) Update(id uuid.UUID, req *request.UpdateSprintRequest, userID uuid.UUID) (*response.SprintResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Track changes
	changes := make(map[string][2]interface{})

	if req.Name != "" && req.Name != sprint.Name {
		changes["name"] = [2]interface{}{sprint.Name, req.Name}
		sprint.Name = strings.TrimSpace(req.Name)
	}

	if req.Goal != "" {
		oldGoal := ""
		if sprint.Goal != nil {
			oldGoal = *sprint.Goal
		}
		if req.Goal != oldGoal {
			changes["goal"] = [2]interface{}{oldGoal, req.Goal}
			goal := strings.TrimSpace(req.Goal)
			sprint.Goal = &goal
		}
	}

	if !req.StartDate.IsZero() && req.StartDate != sprint.StartDate {
		changes["start_date"] = [2]interface{}{sprint.StartDate, req.StartDate}
		sprint.StartDate = req.StartDate
	}

	if !req.EndDate.IsZero() && req.EndDate != sprint.EndDate {
		changes["end_date"] = [2]interface{}{sprint.EndDate, req.EndDate}
		sprint.EndDate = req.EndDate
	}

	// Validate date range after updates
	if !sprint.EndDate.After(sprint.StartDate) {
		return nil, ErrInvalidDateRange
	}

	if err := s.sprintRepo.Update(sprint); err != nil {
		return nil, err
	}

	// Record history for changes
	if len(changes) > 0 {
		oldVals, _ := json.Marshal(changes)
		s.recordSprintHistory(id, userID, nil, constants.SprintActionCreated, nil, datatypes.JSON(oldVals))
	}

	// Fetch updated sprint
	updated, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToSprintResponse(updated), nil
}

func (s *sprintService) Delete(id uuid.UUID) error {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return err
	}
	if sprint == nil {
		return ErrSprintNotFound
	}

	return s.sprintRepo.Delete(id)
}

func (s *sprintService) Start(id uuid.UUID, userID uuid.UUID) (*response.SprintResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Check if sprint is in planning status
	if sprint.Status != constants.SprintStatusPlanning {
		return nil, ErrSprintNotPlanning
	}

	// Check if there's already an active sprint in this project
	activeSprint, err := s.sprintRepo.GetActive(sprint.ProjectID)
	if err != nil {
		return nil, err
	}
	if activeSprint != nil {
		return nil, ErrSprintAlreadyActive
	}

	// Update status
	oldStatus := sprint.Status
	if err := s.sprintRepo.UpdateStatus(id, constants.SprintStatusActive); err != nil {
		return nil, err
	}

	// Record history
	oldVal, _ := json.Marshal(oldStatus)
	newVal, _ := json.Marshal(constants.SprintStatusActive)
	s.recordSprintHistory(id, userID, nil, constants.SprintActionStarted, datatypes.JSON(oldVal), datatypes.JSON(newVal))

	// Fetch updated sprint
	updated, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToSprintResponse(updated), nil
}

func (s *sprintService) Complete(id uuid.UUID, userID uuid.UUID) (*response.SprintResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Check if sprint is active
	if sprint.Status != constants.SprintStatusActive {
		return nil, ErrSprintNotActive
	}

	// Calculate velocity
	velocity, err := s.sprintRepo.CalculateVelocity(id)
	if err != nil {
		return nil, err
	}

	// Update sprint
	sprint.Status = constants.SprintStatusCompleted
	sprint.Velocity = &velocity

	if err := s.sprintRepo.Update(sprint); err != nil {
		return nil, err
	}

	// Record history
	oldVal, _ := json.Marshal(constants.SprintStatusActive)
	newVal, _ := json.Marshal(map[string]interface{}{
		"status":   constants.SprintStatusCompleted,
		"velocity": velocity,
	})
	s.recordSprintHistory(id, userID, nil, constants.SprintActionCompleted, datatypes.JSON(oldVal), datatypes.JSON(newVal))

	// Fetch updated sprint
	updated, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToSprintResponse(updated), nil
}

func (s *sprintService) Cancel(id uuid.UUID, userID uuid.UUID) (*response.SprintResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Check if sprint is active or planning
	if sprint.Status != constants.SprintStatusActive && sprint.Status != constants.SprintStatusPlanning {
		return nil, ErrSprintNotActive
	}

	oldStatus := sprint.Status

	// Update status
	if err := s.sprintRepo.UpdateStatus(id, constants.SprintStatusCancelled); err != nil {
		return nil, err
	}

	// Record history
	oldVal, _ := json.Marshal(oldStatus)
	newVal, _ := json.Marshal(constants.SprintStatusCancelled)
	s.recordSprintHistory(id, userID, nil, constants.SprintActionCancelled, datatypes.JSON(oldVal), datatypes.JSON(newVal))

	// Fetch updated sprint
	updated, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToSprintResponse(updated), nil
}

func (s *sprintService) AddItem(sprintID uuid.UUID, itemID uuid.UUID, userID uuid.UUID) (*response.SprintWithItemsResponse, error) {
	// Check if sprint exists
	sprint, err := s.sprintRepo.GetByID(sprintID)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Check if item exists
	item, err := s.backlogRepo.GetByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	// Check if item is already in this sprint
	if item.SprintID != nil && *item.SprintID == sprintID {
		return nil, ErrItemAlreadyInSprint
	}

	// Update item's sprint
	oldSprintID := item.SprintID
	item.SprintID = &sprintID
	if err := s.backlogRepo.Update(item); err != nil {
		return nil, err
	}

	// Record sprint history (ItemAdded)
	itemVal, _ := json.Marshal(map[string]interface{}{
		"item_id":    itemID,
		"item_title": item.Title,
	})
	s.recordSprintHistory(sprintID, userID, &itemID, constants.SprintActionItemAdded, nil, datatypes.JSON(itemVal))

	// Record item history (SprintAssigned)
	oldVal, _ := json.Marshal(oldSprintID)
	newVal, _ := json.Marshal(sprintID)
	s.recordItemHistory(itemID, userID, constants.ItemActionSprintAssigned, "sprint_id", datatypes.JSON(oldVal), datatypes.JSON(newVal))

	// Fetch sprint with items
	items, err := s.sprintRepo.GetItemsBySprintID(sprintID)
	if err != nil {
		return nil, err
	}

	return response.ToSprintWithItemsResponse(sprint, items), nil
}

func (s *sprintService) RemoveItem(sprintID uuid.UUID, itemID uuid.UUID, userID uuid.UUID) (*response.SprintWithItemsResponse, error) {
	// Check if sprint exists
	sprint, err := s.sprintRepo.GetByID(sprintID)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Check if item exists
	item, err := s.backlogRepo.GetByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrBacklogItemNotFound
	}

	// Check if item is in this sprint
	if item.SprintID == nil || *item.SprintID != sprintID {
		return nil, ErrItemNotInSprint
	}

	// Update item's sprint
	oldSprintID := item.SprintID
	item.SprintID = nil
	if err := s.backlogRepo.Update(item); err != nil {
		return nil, err
	}

	// Record sprint history (ItemRemoved)
	itemVal, _ := json.Marshal(map[string]interface{}{
		"item_id":    itemID,
		"item_title": item.Title,
	})
	s.recordSprintHistory(sprintID, userID, &itemID, constants.SprintActionItemRemoved, datatypes.JSON(itemVal), nil)

	// Record item history (SprintRemoved)
	oldVal, _ := json.Marshal(oldSprintID)
	s.recordItemHistory(itemID, userID, constants.ItemActionSprintRemoved, "sprint_id", datatypes.JSON(oldVal), nil)

	// Fetch sprint with items
	items, err := s.sprintRepo.GetItemsBySprintID(sprintID)
	if err != nil {
		return nil, err
	}

	return response.ToSprintWithItemsResponse(sprint, items), nil
}

func (s *sprintService) GetHistory(id uuid.UUID) ([]response.SprintHistoryResponse, error) {
	// Check if sprint exists
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	histories, err := s.sprintHistoryRepo.GetBySprintID(id)
	if err != nil {
		return nil, err
	}

	return response.ToSprintHistoryListResponse(histories), nil
}

func (s *sprintService) GetReport(id uuid.UUID) (*response.SprintReportResponse, error) {
	sprint, err := s.sprintRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, ErrSprintNotFound
	}

	// Get items in sprint
	items, err := s.sprintRepo.GetItemsBySprintID(id)
	if err != nil {
		return nil, err
	}

	// Calculate stats
	totalItems := len(items)
	completedItems := 0
	totalStoryPoints := 0
	completedStoryPoints := 0

	for _, item := range items {
		if item.StoryPoints != nil {
			totalStoryPoints += *item.StoryPoints
		}
		if item.Status == constants.ItemStatusDone {
			completedItems++
			if item.StoryPoints != nil {
				completedStoryPoints += *item.StoryPoints
			}
		}
	}

	// Calculate completion percentage
	var completionPercentage float64
	if totalItems > 0 {
		completionPercentage = float64(completedItems) / float64(totalItems) * 100
	}

	velocity := completedStoryPoints
	if sprint.Velocity != nil {
		velocity = *sprint.Velocity
	}

	return &response.SprintReportResponse{
		Sprint:               *response.ToSprintResponse(sprint),
		TotalItems:           totalItems,
		CompletedItems:       completedItems,
		TotalStoryPoints:     totalStoryPoints,
		CompletedStoryPoints: completedStoryPoints,
		Velocity:             velocity,
		CompletionPercentage: completionPercentage,
	}, nil
}

// recordSprintHistory is a helper function to record sprint history
func (s *sprintService) recordSprintHistory(sprintID, userID uuid.UUID, itemID *uuid.UUID, action constants.SprintAction, oldValue, newValue datatypes.JSON) {
	history := &models.SprintHistory{
		SprintID:  sprintID,
		UserID:    userID,
		ItemID:    itemID,
		Action:    action,
		OldValue:  oldValue,
		NewValue:  newValue,
		Timestamp: time.Now(),
	}
	s.sprintHistoryRepo.Create(history)
}

// recordItemHistory is a helper function to record item history
func (s *sprintService) recordItemHistory(itemID, userID uuid.UUID, action constants.ItemAction, field string, oldValue, newValue datatypes.JSON) {
	history := &models.ItemHistory{
		ItemID:       itemID,
		UserID:       userID,
		Action:       action,
		FieldChanged: &field,
		OldValue:     oldValue,
		NewValue:     newValue,
		Timestamp:    time.Now(),
	}
	s.itemHistoryRepo.Create(history)
}
