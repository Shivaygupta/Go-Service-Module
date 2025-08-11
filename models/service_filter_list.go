package models

import (
	"strings"

	apperr "kong-go-assignment/errors"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
	DefaultSort  = "created_at"
	DefaultOrder = "desc"
)

var AllowedSortFields = map[string]bool{
	"id":         true,
	"name":       true,
	"created_at": true,
	"updated_at": true,
}

type ServiceListFilter struct {
	Name   string
	SortBy string
	Order  string
	Page   int
	Limit  int
}

func (f *ServiceListFilter) ValidateAndNormalize() error {
	f.Name = strings.TrimSpace(f.Name)
	if f.SortBy != "" && !AllowedSortFields[f.SortBy] {
		return apperr.New(
			apperr.ErrInvalidInput.StatusCode,
			apperr.ErrInvalidInput.Message,
			"Invalid value provided in SortBy field",
		)
	}
	if f.SortBy == "" {
		f.SortBy = DefaultSort
	}

	order := strings.ToLower(f.Order)
	if order != "" && order != "asc" && order != "desc" {
		return apperr.New(
			apperr.ErrInvalidInput.StatusCode,
			apperr.ErrInvalidInput.Message,
			"Invalid value provided in Order field",
		)
	}
	if order == "" {
		f.Order = DefaultOrder
	} else {
		f.Order = order
	}

	if f.Page < 0 {
		return apperr.New(
			apperr.ErrInvalidInput.StatusCode,
			apperr.ErrInvalidInput.Message,
			"Invalid value provided in Page field",
		)
	}
	if f.Page == 0 {
		f.Page = 1
	}

	if f.Limit < 0 {
		return apperr.New(
			apperr.ErrInvalidInput.StatusCode,
			apperr.ErrInvalidInput.Message,
			"Invalid value provided in Limit field",
		)
	}

	if f.Limit == 0 {
		f.Limit = DefaultLimit
	} else if f.Limit > MaxLimit {
		f.Limit = MaxLimit
	}
	return nil
}
