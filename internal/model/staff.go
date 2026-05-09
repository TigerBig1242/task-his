package model

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type StaffRole string

const (
	RoleAdmin StaffRole = "admin"
	RoleStaff StaffRole = "staff"
)

type Staff struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	HospitalID   uuid.UUID  `gorm:"type:uuid;not null" json:"hospital_id"`
	Username     string     `gorm:"type:varchar(50);unique;not null" json:"username"`
	PasswordHash string     `gorm:"type:text;not null" json:"-"`
	Role         StaffRole  `gorm:"type:varchar(20);default:'staff'" json:"role"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	//DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // สำหรับ Soft Delete

	// Relationship: Staff belongs to Hospital
	// GORM จะใช้ HospitalID เป็น Foreign Key โดยอัตโนมัติ
	Hospital Hospital `gorm:"foreignKey:HospitalID" json:"hospital,omitempty"`
}
