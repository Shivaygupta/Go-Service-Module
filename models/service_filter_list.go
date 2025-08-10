package models

import "strings"

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

func (f *ServiceListFilter) Normalize() {

	f.Name = strings.TrimSpace(f.Name)

	if !AllowedSortFields[f.SortBy] {
		f.SortBy = DefaultSort
	}

	order := strings.ToLower(f.Order)
	if order != "asc" && order != "desc" {
		f.Order = DefaultOrder
	} else {
		f.Order = order
	}

	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = DefaultLimit
	} else if f.Limit > MaxLimit {
		f.Limit = MaxLimit
	}
}
