package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/dto/response"
	"sprint-backlog/internal/models"
	"sprint-backlog/internal/repository"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectKeyExists     = errors.New("project key already exists")
	ErrInvalidProjectKey    = errors.New("project key must be uppercase alphanumeric")
)

type ProjectService interface {
	Create(req *request.CreateProjectRequest, userID uuid.UUID) (*response.ProjectResponse, error)
	GetByID(id uuid.UUID) (*response.ProjectResponse, error)
	GetByKey(key string) (*response.ProjectResponse, error)
	GetAll() ([]response.ProjectResponse, error)
	GetAllWithPagination(page, limit int) (*response.ProjectListResponse, error)
	Update(id uuid.UUID, req *request.UpdateProjectRequest) (*response.ProjectResponse, error)
	Delete(id uuid.UUID) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) Create(req *request.CreateProjectRequest, userID uuid.UUID) (*response.ProjectResponse, error) {
	// Normalize key to uppercase
	key := strings.ToUpper(strings.TrimSpace(req.Key))

	// Check if project key already exists
	existing, err := s.projectRepo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrProjectKeyExists
	}

	// Create project
	project := &models.Project{
		Name:        strings.TrimSpace(req.Name),
		Key:         key,
		CreatedByID: userID,
	}

	// Set description if provided
	if req.Description != "" {
		desc := strings.TrimSpace(req.Description)
		project.Description = &desc
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Fetch the created project with relations
	created, err := s.projectRepo.GetByID(project.ID)
	if err != nil {
		return nil, err
	}

	return response.ToProjectResponse(created), nil
}

func (s *projectService) GetByID(id uuid.UUID) (*response.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	return response.ToProjectResponse(project), nil
}

func (s *projectService) GetByKey(key string) (*response.ProjectResponse, error) {
	project, err := s.projectRepo.GetByKey(strings.ToUpper(key))
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	return response.ToProjectResponse(project), nil
}

func (s *projectService) GetAll() ([]response.ProjectResponse, error) {
	projects, err := s.projectRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]response.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *response.ToProjectResponse(&project)
	}

	return responses, nil
}

func (s *projectService) GetAllWithPagination(page, limit int) (*response.ProjectListResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	projects, total, err := s.projectRepo.GetAllWithPagination(page, limit)
	if err != nil {
		return nil, err
	}

	return response.ToProjectListResponse(projects, total, page, limit), nil
}

func (s *projectService) Update(id uuid.UUID, req *request.UpdateProjectRequest) (*response.ProjectResponse, error) {
	// Get existing project
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	// Update fields if provided
	if req.Name != "" {
		project.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		desc := strings.TrimSpace(req.Description)
		project.Description = &desc
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	// Fetch updated project with relations
	updated, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return response.ToProjectResponse(updated), nil
}

func (s *projectService) Delete(id uuid.UUID) error {
	// Check if project exists
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	return s.projectRepo.Delete(id)
}
