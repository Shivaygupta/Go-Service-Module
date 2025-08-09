package models

import "time"

type ServiceVersion struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ServiceID uint      `json:"service_id"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}
