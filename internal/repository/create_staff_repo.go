package repository

import (
	"context"
	// "database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bigthamm/task-his/internal/domain"
	"github.com/bigthamm/task-his/internal/model"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffDomain {
	return &staffRepository{db: db}
}

func (repo staffRepository) FindHospitalByCode(ctx context.Context, code string) (*model.Hospital, error) {
	var hospital model.Hospital

	result := repo.db.WithContext(ctx).
		Where("code = ? AND is_active = ?", code, true).
		First(&hospital)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("staffRepository.FindHospitalByCode: %w", result.Error)
	}

	return &hospital, nil
}

func (repo staffRepository) ExistsByUsernameAndHospital(ctx context.Context, username string, hospitalID string) (bool, error) {
	var count int64

	result := repo.db.WithContext(ctx).
		Model(&model.Staff{}).
		Where("username = ? AND hospital_id = ?", username, hospitalID).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("staffRepository.ExistsByUsernameAndHospital: %w", result.Error)
	}

	return count > 0, nil
}

func (repo staffRepository) CreateStaff(ctx context.Context, staff *model.Staff) error {
	result := repo.db.WithContext(ctx).Create(staff)

	if result.Error != nil {
		return fmt.Errorf("staffRepository.CreateStaff: %w", result.Error)
	}

	return nil
}

func (repo staffRepository) FindByUsernameAndHospitalCode(ctx context.Context, username, hospitalCode string) (*model.Staff, error) {
	var staff model.Staff

	result := repo.db.WithContext(ctx).
		Joins("INNER JOIN hospitals h ON h.id = staffs.hospital_id").
		Where("staffs.username = ?", username).
		Where("h.code = ?", hospitalCode).
		Where("staffs.is_active = ?", true).
		Where("h.is_active = ?", true).
		First(&staff)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("staffRepository.FindByUsernameAndHospitalCode: %w", result.Error)
	}

	return &staff, nil
}

func (repo *staffRepository) UpdateLastLogin(ctx context.Context, staffID string, loginAt time.Time) error {
	result := repo.db.WithContext(ctx).
		Model(&model.Staff{}).
		Where("id = ?", staffID).
		Update("last_login_at", loginAt)

	if result.Error != nil {
		return fmt.Errorf("staffRepository.UpdateLastLogin: %w", result.Error)
	}

	return nil
}
