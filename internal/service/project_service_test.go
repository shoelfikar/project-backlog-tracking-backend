package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/models"
)

// MockProjectRepository is a mock implementation of ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByKey(key string) (*models.Project, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetAll() ([]models.Project, error) {
	args := m.Called()
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetAllWithPagination(page, limit int) ([]models.Project, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectRepository) Update(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProjectService_Create(t *testing.T) {
	t.Run("should create project successfully", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		userID := uuid.New()
		projectID := uuid.New()

		req := &request.CreateProjectRequest{
			Name:        "Test Project",
			Key:         "TP",
			Description: "Test description",
		}

		// GetByKey returns nil (no existing project)
		mockRepo.On("GetByKey", "TP").Return(nil, nil)

		// Create succeeds
		mockRepo.On("Create", mock.AnythingOfType("*models.Project")).Return(nil).Run(func(args mock.Arguments) {
			project := args.Get(0).(*models.Project)
			project.ID = projectID
		})

		// GetByID returns the created project
		createdProject := &models.Project{
			ID:          projectID,
			Name:        "Test Project",
			Key:         "TP",
			CreatedByID: userID,
			CreatedBy: models.User{
				ID:   userID,
				Name: "Test User",
			},
		}
		desc := "Test description"
		createdProject.Description = &desc
		mockRepo.On("GetByID", projectID).Return(createdProject, nil)

		result, err := service.Create(req, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Project", result.Name)
		assert.Equal(t, "TP", result.Key)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when project key already exists", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		userID := uuid.New()
		existingProject := &models.Project{
			ID:   uuid.New(),
			Name: "Existing Project",
			Key:  "EXIST",
		}

		req := &request.CreateProjectRequest{
			Name: "New Project",
			Key:  "exist", // lowercase should be normalized
		}

		mockRepo.On("GetByKey", "EXIST").Return(existingProject, nil)

		result, err := service.Create(req, userID)

		assert.Error(t, err)
		assert.Equal(t, ErrProjectKeyExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should normalize project key to uppercase", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		userID := uuid.New()
		projectID := uuid.New()

		req := &request.CreateProjectRequest{
			Name: "Test Project",
			Key:  "  tp  ", // with spaces and lowercase
		}

		mockRepo.On("GetByKey", "TP").Return(nil, nil)
		mockRepo.On("Create", mock.MatchedBy(func(p *models.Project) bool {
			return p.Key == "TP"
		})).Return(nil).Run(func(args mock.Arguments) {
			args.Get(0).(*models.Project).ID = projectID
		})

		createdProject := &models.Project{
			ID:          projectID,
			Name:        "Test Project",
			Key:         "TP",
			CreatedByID: userID,
		}
		mockRepo.On("GetByID", projectID).Return(createdProject, nil)

		result, err := service.Create(req, userID)

		assert.NoError(t, err)
		assert.Equal(t, "TP", result.Key)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_GetByID(t *testing.T) {
	t.Run("should return project by ID", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		project := &models.Project{
			ID:   projectID,
			Name: "Test Project",
			Key:  "TP",
		}

		mockRepo.On("GetByID", projectID).Return(project, nil)

		result, err := service.GetByID(projectID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, projectID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when project not found", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		mockRepo.On("GetByID", projectID).Return(nil, nil)

		result, err := service.GetByID(projectID)

		assert.Error(t, err)
		assert.Equal(t, ErrProjectNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error on repository error", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		mockRepo.On("GetByID", projectID).Return(nil, errors.New("db error"))

		result, err := service.GetByID(projectID)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_GetAll(t *testing.T) {
	t.Run("should return all projects", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projects := []models.Project{
			{ID: uuid.New(), Name: "Project 1", Key: "P1"},
			{ID: uuid.New(), Name: "Project 2", Key: "P2"},
		}

		mockRepo.On("GetAll").Return(projects, nil)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return empty slice when no projects", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		mockRepo.On("GetAll").Return([]models.Project{}, nil)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_GetAllWithPagination(t *testing.T) {
	t.Run("should return paginated projects", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projects := []models.Project{
			{ID: uuid.New(), Name: "Project 1", Key: "P1"},
		}

		mockRepo.On("GetAllWithPagination", 1, 10).Return(projects, int64(15), nil)

		result, err := service.GetAllWithPagination(1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Projects, 1)
		assert.Equal(t, int64(15), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should use default values for invalid page and limit", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		mockRepo.On("GetAllWithPagination", 1, 10).Return([]models.Project{}, int64(0), nil)

		result, err := service.GetAllWithPagination(0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.Limit)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should cap limit at 100", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		mockRepo.On("GetAllWithPagination", 1, 100).Return([]models.Project{}, int64(0), nil)

		result, err := service.GetAllWithPagination(1, 500)

		assert.NoError(t, err)
		assert.Equal(t, 100, result.Limit)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_Update(t *testing.T) {
	t.Run("should update project successfully", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		existingProject := &models.Project{
			ID:   projectID,
			Name: "Original Name",
			Key:  "ON",
		}

		req := &request.UpdateProjectRequest{
			Name:        "Updated Name",
			Description: "Updated description",
		}

		mockRepo.On("GetByID", projectID).Return(existingProject, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Project")).Return(nil)

		updatedProject := &models.Project{
			ID:   projectID,
			Name: "Updated Name",
			Key:  "ON",
		}
		desc := "Updated description"
		updatedProject.Description = &desc
		mockRepo.On("GetByID", projectID).Return(updatedProject, nil).Once()

		result, err := service.Update(projectID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Name", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when project not found", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		req := &request.UpdateProjectRequest{
			Name: "Updated Name",
		}

		mockRepo.On("GetByID", projectID).Return(nil, nil)

		result, err := service.Update(projectID, req)

		assert.Error(t, err)
		assert.Equal(t, ErrProjectNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProjectService_Delete(t *testing.T) {
	t.Run("should delete project successfully", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		project := &models.Project{
			ID:   projectID,
			Name: "Test Project",
		}

		mockRepo.On("GetByID", projectID).Return(project, nil)
		mockRepo.On("Delete", projectID).Return(nil)

		err := service.Delete(projectID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when project not found", func(t *testing.T) {
		mockRepo := new(MockProjectRepository)
		service := NewProjectService(mockRepo)

		projectID := uuid.New()
		mockRepo.On("GetByID", projectID).Return(nil, nil)

		err := service.Delete(projectID)

		assert.Error(t, err)
		assert.Equal(t, ErrProjectNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
