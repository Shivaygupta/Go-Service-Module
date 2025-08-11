package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"

	apperr "kong-go-assignment/errors"
	"kong-go-assignment/models"
)

var ErrNotFound = errors.New("not found")

type ServiceRepository interface {
	FindAll(ctx context.Context, filter models.ServiceListFilter) (models.PaginatedServices, error)
	FindByID(ctx context.Context, id uint) (*models.Service, error)
	CreateService(ctx context.Context, svc *models.Service) error
	DeleteServiceByID(ctx context.Context, id uint) error
}

type serviceRepo struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepo{db: db}
}

func buildOrderClause(sortBy, order string) string {
	return sortBy + " " + order
}

func calculateOffset(page, limit int) int {
	return (page - 1) * limit
}

func (r *serviceRepo) FindAll(ctx context.Context, filter models.ServiceListFilter) (models.PaginatedServices, error) {
	var out []models.Service
	dbq := r.db.WithContext(ctx).Model(&models.Service{}).Preload("Versions")

	if filter.Name != "" {
		like := "%" + strings.ToLower(filter.Name) + "%"
		dbq = dbq.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", like, like)
	}

	var total int64
	if err := dbq.Count(&total).Error; err != nil {
		return models.PaginatedServices{}, apperr.New(
			apperr.ErrInternal.StatusCode,
			apperr.ErrInternal.Message,
			"serviceRepo.FindAll count error: "+err.Error(),
		)
	}

	dbq = dbq.Order(buildOrderClause(filter.SortBy, filter.Order)).
		Limit(filter.Limit).
		Offset(calculateOffset(filter.Page, filter.Limit))

	if err := dbq.Find(&out).Error; err != nil {
		return models.PaginatedServices{}, apperr.New(
			apperr.ErrInternal.StatusCode,
			apperr.ErrInternal.Message,
			"serviceRepo.FindAll query error: "+err.Error(),
		)
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return models.PaginatedServices{
		Data:       out,
		TotalCount: total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *serviceRepo) FindByID(ctx context.Context, id uint) (*models.Service, error) {
	var svc models.Service
	if err := r.db.WithContext(ctx).Preload("Versions").First(&svc, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(
				apperr.ErrInternal.StatusCode,
				apperr.ErrInternal.Message,
				fmt.Sprintf("Service with id: %d not found", id))
		}
		return nil, err
	}
	return &svc, nil
}

func (r *serviceRepo) CreateService(ctx context.Context, svc *models.Service) error {
	var existing models.Service
	err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", svc.Name).First(&existing).Error
	if err == nil {
		return apperr.New(
			apperr.ErrConflict.StatusCode,
			"Service with the same name already exists",
			"duplicate service name: "+svc.Name,
		)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.New(
			apperr.ErrInternal.StatusCode,
			apperr.ErrInternal.Message,
			"serviceRepo.Create duplicate check error: "+err.Error(),
		)
	}

	if len(svc.Versions) > 0 {
		for i := range svc.Versions {
			svc.Versions[i].ID = svc.ID
		}
	}

	if err := r.db.WithContext(ctx).Create(svc).Error; err != nil {
		return apperr.New(
			apperr.ErrInternal.StatusCode,
			apperr.ErrInternal.Message,
			"serviceRepo.Create insert error: "+err.Error(),
		)
	}

	return nil
}

func (r *serviceRepo) DeleteServiceByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&models.Service{})

		if result.Error != nil {
			return apperr.New(
				apperr.ErrInternal.StatusCode,
				apperr.ErrInternal.Message,
				"serviceRepo.DeleteServiceByID delete error: "+result.Error.Error(),
			)
		}

		if result.RowsAffected == 0 {
			return apperr.ErrServiceNotFound
		}
		log.Printf("Service with ID %d deleted successfully", id)
		return nil
	})
}
