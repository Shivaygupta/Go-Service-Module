package repositories

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	apperr "kong-go-assignment/errors"
	"kong-go-assignment/models"
)

var ErrNotFound = errors.New("not found")

type ServiceRepository interface {
	FindAll(ctx context.Context, filter models.ServiceListFilter) (models.PaginatedServices, error)
	FindByID(ctx context.Context, id uint) (*models.Service, error)
	Create(ctx context.Context, svc *models.Service) error
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

	filter.Normalize()

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
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &svc, nil
}

func (r *serviceRepo) Create(ctx context.Context, svc *models.Service) error {
	return r.db.WithContext(ctx).Create(svc).Error
}
