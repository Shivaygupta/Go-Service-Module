package models

import "time"

type Service struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Versions    []*ServiceVersion `json:"versions" gorm:"foreignKey:ServiceID"`
}
