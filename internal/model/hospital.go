package model

import (
	"time"

	"github.com/google/uuid"
)

type Hospital struct {
	// ID         uuid.UUID `json:"id"          db:"id"`
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `json:"name"        db:"name"`
	Code       string    `json:"code"        db:"code"`
	HISBaseURL string    `json:"-"           db:"his_base_url"` // ไม่ expose ออก API
	HISAPIKey  string    `json:"-"           db:"his_api_key"`  // ไม่ expose ออก API
	Address    string    `json:"address"     db:"address"`
	IsActive   bool      `json:"is_active"   db:"is_active"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
}
