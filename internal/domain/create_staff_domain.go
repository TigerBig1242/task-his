package domain

import (
	"context"
	"time"

	"github.com/bigthamm/task-his/internal/model"
)

type StaffDomain interface {
	CreateStaff(ctx context.Context, staff *model.Staff) error
	FindHospitalByCode(ctx context.Context, code string) (*model.Hospital, error)
	ExistsByUsernameAndHospital(ctx context.Context, username string, hospitalID string) (bool, error)
	FindByUsernameAndHospitalCode(ctx context.Context, username, hospitalCode string) (*model.Staff, error)
	UpdateLastLogin(ctx context.Context, staffID string, loginAt time.Time) error
}
