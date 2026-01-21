package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/dto/response"
	"sprint-backlog/internal/service"
)

// MockProjectService is a mock implementation of ProjectService
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) Create(req *request.CreateProjectRequest, userID uuid.UUID) (*response.ProjectResponse, error) {
	args := m.Called(req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) GetByID(id uuid.UUID) (*response.ProjectResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) GetByKey(key string) (*response.ProjectResponse, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) GetAll() ([]response.ProjectResponse, error) {
	args := m.Called()
	return args.Get(0).([]response.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) GetAllWithPagination(page, limit int) (*response.ProjectListResponse, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ProjectListResponse), args.Error(1)
}

func (m *MockProjectService) Update(id uuid.UUID, req *request.UpdateProjectRequest) (*response.ProjectResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ProjectResponse), args.Error(1)
}

func (m *MockProjectService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestProjectHandler_GetAll(t *testing.T) {
	t.Run("should return projects with pagination", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.GET("/projects", handler.GetAll)

		projectList := &response.ProjectListResponse{
			Projects: []response.ProjectResponse{
				{ID: uuid.New(), Name: "Project 1", Key: "P1"},
			},
			Total:      1,
			Page:       1,
			Limit:      10,
			TotalPages: 1,
		}
		mockService.On("GetAllWithPagination", 1, 10).Return(projectList, nil)

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should use custom page and limit", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.GET("/projects", handler.GetAll)

		projectList := &response.ProjectListResponse{
			Projects:   []response.ProjectResponse{},
			Total:      0,
			Page:       2,
			Limit:      20,
			TotalPages: 0,
		}
		mockService.On("GetAllWithPagination", 2, 20).Return(projectList, nil)

		req := httptest.NewRequest(http.MethodGet, "/projects?page=2&limit=20", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProjectHandler_GetByID(t *testing.T) {
	t.Run("should return project by ID", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.GET("/projects/:id", handler.GetByID)

		projectID := uuid.New()
		project := &response.ProjectResponse{
			ID:   projectID,
			Name: "Test Project",
			Key:  "TP",
		}
		mockService.On("GetByID", projectID).Return(project, nil)

		req := httptest.NewRequest(http.MethodGet, "/projects/"+projectID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.GET("/projects/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/projects/invalid-uuid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 404 when project not found", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.GET("/projects/:id", handler.GetByID)

		projectID := uuid.New()
		mockService.On("GetByID", projectID).Return(nil, service.ErrProjectNotFound)

		req := httptest.NewRequest(http.MethodGet, "/projects/"+projectID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProjectHandler_Create(t *testing.T) {
	t.Run("should create project successfully", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		// Add middleware to set user ID in context
		router.POST("/projects", func(c *gin.Context) {
			c.Set("user_id", uuid.New())
			handler.Create(c)
		})

		createReq := request.CreateProjectRequest{
			Name: "New Project",
			Key:  "NP",
		}

		project := &response.ProjectResponse{
			ID:   uuid.New(),
			Name: "New Project",
			Key:  "NP",
		}
		mockService.On("Create", mock.AnythingOfType("*request.CreateProjectRequest"), mock.AnythingOfType("uuid.UUID")).Return(project, nil)

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid request body", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.POST("/projects", func(c *gin.Context) {
			c.Set("user_id", uuid.New())
			handler.Create(c)
		})

		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 401 when user not authenticated", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.POST("/projects", handler.Create)

		createReq := request.CreateProjectRequest{
			Name: "New Project",
			Key:  "NP",
		}

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return 409 when project key exists", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.POST("/projects", func(c *gin.Context) {
			c.Set("user_id", uuid.New())
			handler.Create(c)
		})

		createReq := request.CreateProjectRequest{
			Name: "New Project",
			Key:  "EXIST",
		}

		mockService.On("Create", mock.AnythingOfType("*request.CreateProjectRequest"), mock.AnythingOfType("uuid.UUID")).Return(nil, service.ErrProjectKeyExists)

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProjectHandler_Update(t *testing.T) {
	t.Run("should update project successfully", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.PUT("/projects/:id", handler.Update)

		projectID := uuid.New()
		updateReq := request.UpdateProjectRequest{
			Name: "Updated Project",
		}

		project := &response.ProjectResponse{
			ID:   projectID,
			Name: "Updated Project",
			Key:  "UP",
		}
		mockService.On("Update", projectID, mock.AnythingOfType("*request.UpdateProjectRequest")).Return(project, nil)

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/projects/"+projectID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 when project not found", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.PUT("/projects/:id", handler.Update)

		projectID := uuid.New()
		updateReq := request.UpdateProjectRequest{
			Name: "Updated Project",
		}

		mockService.On("Update", projectID, mock.AnythingOfType("*request.UpdateProjectRequest")).Return(nil, service.ErrProjectNotFound)

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/projects/"+projectID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProjectHandler_Delete(t *testing.T) {
	t.Run("should delete project successfully", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.DELETE("/projects/:id", handler.Delete)

		projectID := uuid.New()
		mockService.On("Delete", projectID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/projects/"+projectID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 when project not found", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.DELETE("/projects/:id", handler.Delete)

		projectID := uuid.New()
		mockService.On("Delete", projectID).Return(service.ErrProjectNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/projects/"+projectID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid UUID", func(t *testing.T) {
		mockService := new(MockProjectService)
		handler := NewProjectHandler(mockService)

		router := setupTestRouter()
		router.DELETE("/projects/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/projects/invalid-uuid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
