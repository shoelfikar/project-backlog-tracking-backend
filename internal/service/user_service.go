package service

import (
	"encoding/json"
	"sort"

	"github.com/google/uuid"

	"sprint-backlog/internal/dto/response"
	"sprint-backlog/internal/repository"
)

type UserService interface {
	GetAll() ([]response.UserResponse, error)
	GetByID(id uuid.UUID) (*response.UserResponse, error)
	GetActivities(userID uuid.UUID, limit int) (*response.UserActivitiesResponse, error)
}

type userService struct {
	userRepo          repository.UserRepository
	itemHistoryRepo   repository.ItemHistoryRepository
	sprintHistoryRepo repository.SprintHistoryRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	itemHistoryRepo repository.ItemHistoryRepository,
	sprintHistoryRepo repository.SprintHistoryRepository,
) UserService {
	return &userService{
		userRepo:          userRepo,
		itemHistoryRepo:   itemHistoryRepo,
		sprintHistoryRepo: sprintHistoryRepo,
	}
}

func (s *userService) GetAll() ([]response.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return response.ToUserListResponse(users), nil
}

func (s *userService) GetByID(id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return response.ToUserResponse(user), nil
}

func (s *userService) GetActivities(userID uuid.UUID, limit int) (*response.UserActivitiesResponse, error) {
	// Get item histories for user
	itemHistories, err := s.itemHistoryRepo.GetByUserID(userID, 0) // Get all, we'll limit after merging
	if err != nil {
		return nil, err
	}

	// Get sprint histories for user
	sprintHistories, err := s.sprintHistoryRepo.GetByUserID(userID, 0) // Get all, we'll limit after merging
	if err != nil {
		return nil, err
	}

	// Convert and merge activities
	activities := make([]response.UserActivityResponse, 0, len(itemHistories)+len(sprintHistories))

	// Add item history activities
	for _, h := range itemHistories {
		activity := response.UserActivityResponse{
			ID:           h.ID,
			Type:         "item",
			Action:       string(h.Action),
			Timestamp:    h.Timestamp,
			FieldChanged: h.FieldChanged,
			Comment:      h.Comment,
		}

		// Parse JSON values
		if h.OldValue != nil {
			var oldVal interface{}
			if err := json.Unmarshal(h.OldValue, &oldVal); err == nil {
				activity.OldValue = oldVal
			}
		}
		if h.NewValue != nil {
			var newVal interface{}
			if err := json.Unmarshal(h.NewValue, &newVal); err == nil {
				activity.NewValue = newVal
			}
		}

		// Add item summary if available
		if h.Item.ID != uuid.Nil {
			activity.Item = &response.BacklogItemSummary{
				ID:    h.Item.ID,
				Title: h.Item.Title,
				Type:  string(h.Item.Type),
			}
		}

		activities = append(activities, activity)
	}

	// Add sprint history activities
	for _, h := range sprintHistories {
		activity := response.UserActivityResponse{
			ID:        h.ID,
			Type:      "sprint",
			Action:    string(h.Action),
			Timestamp: h.Timestamp,
		}

		// Parse JSON values
		if h.OldValue != nil {
			var oldVal interface{}
			if err := json.Unmarshal(h.OldValue, &oldVal); err == nil {
				activity.OldValue = oldVal
			}
		}
		if h.NewValue != nil {
			var newVal interface{}
			if err := json.Unmarshal(h.NewValue, &newVal); err == nil {
				activity.NewValue = newVal
			}
		}

		// Add sprint summary if available
		if h.Sprint.ID != uuid.Nil {
			activity.Sprint = &response.SprintSummaryFull{
				ID:   h.Sprint.ID,
				Name: h.Sprint.Name,
			}
		}

		// Add item summary if available (for item-related sprint actions)
		if h.Item != nil && h.Item.ID != uuid.Nil {
			activity.Item = &response.BacklogItemSummary{
				ID:    h.Item.ID,
				Title: h.Item.Title,
				Type:  string(h.Item.Type),
			}
		}

		activities = append(activities, activity)
	}

	// Sort by timestamp descending
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp.After(activities[j].Timestamp)
	})

	// Apply limit
	total := len(activities)
	if limit > 0 && len(activities) > limit {
		activities = activities[:limit]
	}

	return &response.UserActivitiesResponse{
		Activities: activities,
		Total:      total,
		Limit:      limit,
	}, nil
}
