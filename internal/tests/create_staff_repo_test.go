package tests

import (
	"context"

	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/bigthamm/task-his/internal/model"
	"github.com/bigthamm/task-his/internal/repository"
)

// type mockStaffRepo struct {
// 	mock.Mock
// }

type pgError struct {
	Code    string
	Message string
}

func (e *pgError) Error() string { return e.Message }

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // ปิด log ใน test
	})
	require.NoError(t, err)

	t.Cleanup(func() { sqlDB.Close() })

	return db, mock
}

func hospitalRow(h *model.Hospital) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "name", "code", "his_base_url", "his_api_key",
		"address", "is_active", "created_at", "updated_at",
	}).AddRow(
		h.ID, h.Name, h.Code, h.HISBaseURL, h.HISAPIKey,
		h.Address, h.IsActive, h.CreatedAt, h.UpdatedAt,
	)
}

func staffRow(s *model.Staff) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "hospital_id", "username", "password_hash",
		"role", "is_active", "last_login_at", "created_at", "updated_at",
	}).AddRow(
		s.ID, s.HospitalID, s.Username, s.PasswordHash,
		s.Role, s.IsActive, s.LastLoginAt, s.CreatedAt, s.UpdatedAt,
	)
}

func makeTestCreateStaff(hospitalID uuid.UUID) *model.Staff {
	return &model.Staff{
		ID:           uuid.New(),
		HospitalID:   hospitalID,
		Username:     "john_doe",
		PasswordHash: "$2a$10$hashedpassword",
		Role:         model.RoleStaff,
		IsActive:     true,
		LastLoginAt:  &time.Time{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func makeTestHospital() *model.Hospital {
	return &model.Hospital{
		ID:         uuid.New(),
		Name:       "Hospital A",
		Code:       "HOSP-A",
		HISBaseURL: "https://hospital-a.api.co.th",
		HISAPIKey:  "test-key",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// Create staff
func TestStaffRepository_CreateStaff(t *testing.T) {

	t.Run("positive - inserts staff successfully", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := repository.NewStaffRepository(db)
		hospital := makeTestHospital()
		staff := makeTestCreateStaff(hospital.ID)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "staffs"`)).
			WithArgs(
				staff.ID, staff.HospitalID, staff.Username,
				staff.PasswordHash, staff.Role, staff.IsActive,
				staff.LastLoginAt, sqlmock.AnyArg(), sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateStaff(context.Background(), staff)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("negative - db error on insert propagates", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := repository.NewStaffRepository(db)
		staff := makeTestCreateStaff(uuid.New())

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "staffs"`)).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.CreateStaff(context.Background(), staff)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CreateStaff")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("negative - duplicate username constraint violation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := repository.NewStaffRepository(db)
		staff := makeTestCreateStaff(uuid.New())

		duplicateErr := &pgError{Code: "23505", Message: "duplicate key value violates unique constraint"}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "staffs"`)).
			WillReturnError(duplicateErr)
		mock.ExpectRollback()

		err := repo.CreateStaff(context.Background(), staff)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CreateStaff")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Find hospital by code
func TestStaffRepository_FindHospitalByCode(t *testing.T) {

	t.Run("positive - found active hospital", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := repository.NewStaffRepository(db)
		hospital := makeTestHospital()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals"`)).
			WithArgs(hospital.Code, true, 1).
			WillReturnRows(hospitalRow(hospital))

		result, err := repo.FindHospitalByCode(context.Background(), hospital.Code)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, hospital.Code, result.Code)
		assert.Equal(t, hospital.ID, result.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("negative - hospital not found returns nil without error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := repository.NewStaffRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals"`)).
			WithArgs("HOSP-INVALID", true, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.FindHospitalByCode(context.Background(), "HOSP-INVALID")

		assert.NoError(t, err) // ไม่ควร error
		assert.Nil(t, result)  // แต่ควรเป็น nil
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("negative - db error propagates", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := repository.NewStaffRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals"`)).
			WithArgs("HOSP-A", true, 1).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.FindHospitalByCode(context.Background(), "HOSP-A")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "FindHospitalByCode")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
