package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bigthamm/task-his/internal/domain"
	"github.com/bigthamm/task-his/internal/dto"
	"github.com/bigthamm/task-his/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StaffService interface {
	Create(ctx context.Context, req dto.CreateStaffRequest) (*dto.StaffResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type staffService struct {
	staffRepo domain.StaffDomain
}

func NewStaffService(staffRepo domain.StaffDomain) StaffService {
	return &staffService{staffRepo: staffRepo}
}

var (
	ErrInvalidCredentials   = errors.New("invalid username, password, or hospital")
	ErrAccountInactive      = errors.New("account is inactive")
	ErrHospitalNotFound     = errors.New("hospital not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken in this hospital")
)

func (service *staffService) Create(ctx context.Context, req dto.CreateStaffRequest) (*dto.StaffResponse, error) {
	hospital, err := service.staffRepo.FindHospitalByCode(ctx, req.Hospital)
	if err != nil {
		return nil, fmt.Errorf("staffService.Create: Find hospital code not found: %w", err)
	}

	if hospital == nil {
		return nil, ErrHospitalNotFound
	}

	usernameExist, err := service.staffRepo.ExistsByUsernameAndHospital(ctx, req.Username, hospital.ID.String())
	if err != nil {
		return nil, fmt.Errorf("staffService.Create: Username already exist: %w", err)
	}

	if usernameExist {
		return nil, ErrUsernameAlreadyTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("staffService.Create: Password has not hash: %w", err)
	}

	now := time.Now()
	staff := &model.Staff{
		ID:           uuid.New(),
		HospitalID:   hospital.ID,
		Username:     req.Username,
		PasswordHash: string(hash),
		Role:         model.RoleStaff,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	createErr := service.staffRepo.CreateStaff(ctx, staff)
	if createErr != nil {
		return nil, fmt.Errorf("staffService.CreateStaff: Create staff fail: %w", err)
	}

	return &dto.StaffResponse{
		ID:         staff.ID,
		HospitalID: staff.HospitalID,
		Username:   staff.Username,
		Role:       staff.Role,
		IsActive:   staff.IsActive,
		CreatedAt:  staff.CreatedAt,
	}, nil
}

// Login staff by validate credentials
func (service *staffService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	staff, err := service.staffRepo.FindByUsernameAndHospitalCode(ctx, req.Username, req.Hospital)
	if err != nil {
		return nil, fmt.Errorf("staffService.Login: Login require username and hospital: %w", err)
	}

	if staff == nil {
		return nil, ErrInvalidCredentials
	}

	errPasswordNotMatch := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password))
	if errPasswordNotMatch != nil {
		return nil, ErrInvalidCredentials
	}

	return &dto.LoginResponse{
		Username:     req.Username,
		HospitalCode: req.Hospital,
	}, nil
}
