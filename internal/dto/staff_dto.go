package dto

import (
	"time"

	"github.com/bigthamm/task-his/internal/model"
	"github.com/google/uuid"
)

// DTO สำหรับ response — ไม่มี sensitive field
type StaffResponse struct {
	ID         uuid.UUID       `json:"id"`
	HospitalID uuid.UUID       `json:"hospital_id"`
	Username   string          `json:"username"`
	Role       model.StaffRole `json:"role"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
}

// Request body สำหรับ /staff/create
type CreateStaffRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	Hospital string `json:"hospital" binding:"required"` // hospital code เช่น HOSP-A
}

// Request body สำหรับ /staff/login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

// Response หลัง login สำเร็จ
type LoginResponse struct {
	Username     string    `json:"username"`
	HospitalCode string    `json:"hospital_code"`
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
