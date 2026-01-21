package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sprint-backlog/internal/dto/response"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll() ([]response.UserResponse, error) {
	args := m.Called()
	return args.Get(0).([]response.UserResponse), args.Error(1)
}

func (m *MockUserService) GetByID(id uuid.UUID) (*response.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func (m *MockUserService) GetActivities(userID uuid.UUID, limit int) (*response.UserActivitiesResponse, error) {
	args := m.Called(userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.UserActivitiesResponse), args.Error(1)
}

func TestUserHandler_GetAll(t *testing.T) {
	t.Run("should return all users", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users", handler.GetAll)

		users := []response.UserResponse{
			{ID: uuid.New(), Name: "User 1", Email: "user1@example.com"},
			{ID: uuid.New(), Name: "User 2", Email: "user2@example.com"},
		}
		mockService.On("GetAll").Return(users, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return empty slice when no users", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users", handler.GetAll)

		mockService.On("GetAll").Return([]response.UserResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("should return user by ID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users/:id", handler.GetByID)

		userID := uuid.New()
		user := &response.UserResponse{
			ID:    userID,
			Name:  "Test User",
			Email: "test@example.com",
		}
		mockService.On("GetByID", userID).Return(user, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 404 when user not found", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users/:id", handler.GetByID)

		userID := uuid.New()
		mockService.On("GetByID", userID).Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetActivities(t *testing.T) {
	t.Run("should return user activities", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users/:id/activities", handler.GetActivities)

		userID := uuid.New()
		activities := &response.UserActivitiesResponse{
			Activities: []response.UserActivityResponse{},
			Total:      0,
			Limit:      50,
		}
		mockService.On("GetActivities", userID, 50).Return(activities, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/activities", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should use custom limit", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users/:id/activities", handler.GetActivities)

		userID := uuid.New()
		activities := &response.UserActivitiesResponse{
			Activities: []response.UserActivityResponse{},
			Total:      0,
			Limit:      100,
		}
		mockService.On("GetActivities", userID, 100).Return(activities, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/activities?limit=100", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := setupTestRouter()
		router.GET("/users/:id/activities", handler.GetActivities)

		req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid/activities", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
