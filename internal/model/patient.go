package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gender string

const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
)

// Patient — map ตรงกับ table patients ใน DB
type Patient struct {
	// ใช้ type uuid.UUID และตั้งเป็น Primary Key
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	HospitalID   uuid.UUID      `gorm:"type:uuid;not null" json:"hospital_id"`
	PatientHN    string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_hospital_hn" json:"patient_hn"`
	NationalID   *string        `gorm:"type:varchar(20);index" json:"national_id"`
	PassportID   *string        `gorm:"type:varchar(20);index" json:"passport_id"`
	FirstNameTH  string         `gorm:"type:varchar(100);not null" json:"first_name_th"`
	MiddleNameTH *string        `gorm:"type:varchar(100)" json:"middle_name_th"`
	LastNameTH   string         `gorm:"type:varchar(100);not null" json:"last_name_th"`
	FirstNameEN  string         `gorm:"type:varchar(100);not null" json:"first_name_en"`
	MiddleNameEN *string        `gorm:"type:varchar(100)" json:"middle_name_en"`
	LastNameEN   string         `gorm:"type:varchar(100);not null" json:"last_name_en"`
	DateOfBirth  time.Time      `gorm:"type:date;not null" json:"date_of_birth"`
	PhoneNumber  *string        `gorm:"type:varchar(20)" json:"phone_number"`
	Email        *string        `gorm:"type:varchar(100)" json:"email"`
	Gender       Gender         `gorm:"type:varchar(10)" json:"gender"`
	SyncedAt     *time.Time     `json:"synced_at"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // สำหรับ Soft Delete
}
