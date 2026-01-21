package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sprint-backlog/internal/models"
	"sprint-backlog/pkg/constants"
)

// MockUserRepository for user service tests
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	args := m.Called(googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

// MockItemHistoryRepository for user service tests
type MockItemHistoryRepository struct {
	mock.Mock
}

func (m *MockItemHistoryRepository) Create(history *models.ItemHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockItemHistoryRepository) GetByItemID(itemID uuid.UUID) ([]models.ItemHistory, error) {
	args := m.Called(itemID)
	return args.Get(0).([]models.ItemHistory), args.Error(1)
}

func (m *MockItemHistoryRepository) GetByUserID(userID uuid.UUID, limit int) ([]models.ItemHistory, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]models.ItemHistory), args.Error(1)
}

// MockSprintHistoryRepository for user service tests
type MockSprintHistoryRepository struct {
	mock.Mock
}

func (m *MockSprintHistoryRepository) Create(history *models.SprintHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockSprintHistoryRepository) GetBySprintID(sprintID uuid.UUID) ([]models.SprintHistory, error) {
	args := m.Called(sprintID)
	return args.Get(0).([]models.SprintHistory), args.Error(1)
}

func (m *MockSprintHistoryRepository) GetByUserID(userID uuid.UUID, limit int) ([]models.SprintHistory, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]models.SprintHistory), args.Error(1)
}

func (m *MockSprintHistoryRepository) GetAll(limit int) ([]models.SprintHistory, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.SprintHistory), args.Error(1)
}

func TestUserService_GetAll(t *testing.T) {
	t.Run("should return all users", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		users := []models.User{
			{ID: uuid.New(), Name: "User 1", Email: "user1@example.com", GoogleID: "g1"},
			{ID: uuid.New(), Name: "User 2", Email: "user2@example.com", GoogleID: "g2"},
		}

		mockUserRepo.On("GetAll").Return(users, nil)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "User 1", result[0].Name)
		assert.Equal(t, "User 2", result[1].Name)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return empty slice when no users", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		mockUserRepo.On("GetAll").Return([]models.User{}, nil)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("should return user by ID", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		userID := uuid.New()
		user := &models.User{
			ID:       userID,
			Name:     "Test User",
			Email:    "test@example.com",
			GoogleID: "google123",
		}

		mockUserRepo.On("GetByID", userID).Return(user, nil)

		result, err := service.GetByID(userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "Test User", result.Name)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return nil when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		userID := uuid.New()
		mockUserRepo.On("GetByID", userID).Return(nil, nil)

		result, err := service.GetByID(userID)

		assert.NoError(t, err)
		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetActivities(t *testing.T) {
	t.Run("should return merged and sorted activities", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		userID := uuid.New()
		itemID := uuid.New()
		sprintID := uuid.New()

		now := time.Now()
		oldTime := now.Add(-1 * time.Hour)

		itemHistories := []models.ItemHistory{
			{
				ID:        uuid.New(),
				ItemID:    itemID,
				UserID:    userID,
				Action:    constants.ItemActionCreated,
				Timestamp: oldTime,
				Item: models.BacklogItem{
					ID:    itemID,
					Title: "Test Item",
					Type:  constants.ItemTypeTask,
				},
			},
		}

		sprintHistories := []models.SprintHistory{
			{
				ID:        uuid.New(),
				SprintID:  sprintID,
				UserID:    userID,
				Action:    constants.SprintActionStarted,
				Timestamp: now,
				Sprint: models.Sprint{
					ID:   sprintID,
					Name: "Sprint 1",
				},
			},
		}

		mockItemHistoryRepo.On("GetByUserID", userID, 0).Return(itemHistories, nil)
		mockSprintHistoryRepo.On("GetByUserID", userID, 0).Return(sprintHistories, nil)

		result, err := service.GetActivities(userID, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Total)
		assert.Len(t, result.Activities, 2)

		// First activity should be the sprint action (newer)
		assert.Equal(t, "sprint", result.Activities[0].Type)
		assert.Equal(t, string(constants.SprintActionStarted), result.Activities[0].Action)

		// Second activity should be the item action (older)
		assert.Equal(t, "item", result.Activities[1].Type)
		assert.Equal(t, string(constants.ItemActionCreated), result.Activities[1].Action)

		mockItemHistoryRepo.AssertExpectations(t)
		mockSprintHistoryRepo.AssertExpectations(t)
	})

	t.Run("should apply limit to activities", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		userID := uuid.New()

		// Create 5 item histories
		itemHistories := make([]models.ItemHistory, 5)
		for i := 0; i < 5; i++ {
			itemHistories[i] = models.ItemHistory{
				ID:        uuid.New(),
				ItemID:    uuid.New(),
				UserID:    userID,
				Action:    constants.ItemActionCreated,
				Timestamp: time.Now().Add(time.Duration(-i) * time.Hour),
			}
		}

		mockItemHistoryRepo.On("GetByUserID", userID, 0).Return(itemHistories, nil)
		mockSprintHistoryRepo.On("GetByUserID", userID, 0).Return([]models.SprintHistory{}, nil)

		result, err := service.GetActivities(userID, 3)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, result.Total)
		assert.Len(t, result.Activities, 3)
		assert.Equal(t, 3, result.Limit)

		mockItemHistoryRepo.AssertExpectations(t)
		mockSprintHistoryRepo.AssertExpectations(t)
	})

	t.Run("should return empty activities when no history", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockItemHistoryRepo := new(MockItemHistoryRepository)
		mockSprintHistoryRepo := new(MockSprintHistoryRepository)

		service := NewUserService(mockUserRepo, mockItemHistoryRepo, mockSprintHistoryRepo)

		userID := uuid.New()

		mockItemHistoryRepo.On("GetByUserID", userID, 0).Return([]models.ItemHistory{}, nil)
		mockSprintHistoryRepo.On("GetByUserID", userID, 0).Return([]models.SprintHistory{}, nil)

		result, err := service.GetActivities(userID, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Total)
		assert.Empty(t, result.Activities)

		mockItemHistoryRepo.AssertExpectations(t)
		mockSprintHistoryRepo.AssertExpectations(t)
	})
}
