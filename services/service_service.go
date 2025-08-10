package services

import (
	"context"

	"kong-go-assignment/errors"
	"kong-go-assignment/models"
	"kong-go-assignment/repositories"
)

type ServiceService interface {
	ListServices(ctx context.Context, filter models.ServiceListFilter) (models.PaginatedServices, error)
	GetService(ctx context.Context, id uint) (*models.Service, error)
	CreateService(ctx context.Context, svc *models.Service) error
	DeleteService(ctx context.Context, id uint) error
}

type serviceService struct {
	repo repositories.ServiceRepository
}

func NewServiceService(repo repositories.ServiceRepository) ServiceService {
	return &serviceService{repo: repo}
}

func (s *serviceService) ListServices(ctx context.Context, filter models.ServiceListFilter) (models.PaginatedServices, error) {
	filter.Normalize()
	return s.repo.FindAll(ctx, filter)
}

func (s *serviceService) GetService(ctx context.Context, id uint) (*models.Service, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *serviceService) CreateService(ctx context.Context, svc *models.Service) error {
	if svc.Name == "" {
		return errors.New(errors.ErrInvalidInput.StatusCode, "service name is required", "empty service name in CreateService")
	}
	return s.repo.CreateService(ctx, svc)
}

func (s *serviceService) DeleteService(ctx context.Context, id uint) error {

	return s.repo.DeleteServiceByID(ctx, id)
}
