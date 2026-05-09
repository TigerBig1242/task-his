package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/bigthamm/task-his/internal/dto"
	"github.com/bigthamm/task-his/internal/model"
	"github.com/bigthamm/task-his/internal/service"
)

type mockStaffRepo struct{ mock.Mock }

func (m *mockStaffRepo) FindHospitalByCode(ctx context.Context, code string) (*model.Hospital, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Hospital), args.Error(1)
}

func (m *mockStaffRepo) ExistsByUsernameAndHospital(ctx context.Context, username, hospitalID string) (bool, error) {
	args := m.Called(ctx, username, hospitalID)
	return args.Bool(0), args.Error(1)
}

func (m *mockStaffRepo) CreateStaff(ctx context.Context, staff *model.Staff) error {
	return m.Called(ctx, staff).Error(0)
}

func (m *mockStaffRepo) FindByUsernameAndHospitalCode(ctx context.Context, username, hospitalCode string) (*model.Staff, error) {
	args := m.Called(ctx, username, hospitalCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Staff), args.Error(1)
}

func (m *mockStaffRepo) UpdateLastLogin(ctx context.Context, staffID string, loginAt time.Time) error {
	return m.Called(ctx, staffID, loginAt).Error(0)
}

func makeHospital(code string) *model.Hospital {
	return &model.Hospital{
		ID: uuid.New(), Name: "Hospital " + code, Code: code,
		HISBaseURL: "https://hospital-a.api.co.th", HISAPIKey: "key",
		IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func makeHashedStaff(username, password string, hospitalID uuid.UUID) *model.Staff {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return &model.Staff{
		ID: uuid.New(), HospitalID: hospitalID, Username: username,
		PasswordHash: string(hash), Role: model.RoleStaff,
		IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

// Test Create Staff
func TestStaffService_Create(t *testing.T) {
	t.Run("positive - creates staff and returns DTO without sensitive fields", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(hospital, nil)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, "john_doe", hospital.ID.String()).Return(false, nil)
		repo.On("CreateStaff", mock.Anything, mock.MatchedBy(func(s *model.Staff) bool {
			return s.Username == "john_doe" &&
				s.HospitalID == hospital.ID &&
				s.PasswordHash != "password123" && // ต้องไม่เป็น plaintext
				s.Role == model.RoleStaff &&
				s.IsActive == true
		})).Return(nil)

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "john_doe", resp.Username)
		assert.Equal(t, hospital.ID, resp.HospitalID)
		assert.Equal(t, model.RoleStaff, resp.Role)
		assert.True(t, resp.IsActive)
		assert.NotEqual(t, uuid.Nil, resp.ID)
		repo.AssertExpectations(t)
	})

	t.Run("positive - password is bcrypt hashed correctly", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")
		plain := "securePass99"

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(hospital, nil)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
		repo.On("CreateStaff", mock.Anything, mock.MatchedBy(func(s *model.Staff) bool {
			return bcrypt.CompareHashAndPassword([]byte(s.PasswordHash), []byte(plain)) == nil
		})).Return(nil)

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "alice", Password: plain, Hospital: "HOSP-A",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("positive - each staff gets unique uuid", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(hospital, nil).Times(2)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Times(2)
		repo.On("CreateStaff", mock.Anything, mock.Anything).Return(nil).Times(2)

		resp1, _ := svc.Create(context.Background(), dto.CreateStaffRequest{Username: "user1", Password: "pass1234", Hospital: "HOSP-A"})
		resp2, _ := svc.Create(context.Background(), dto.CreateStaffRequest{Username: "user2", Password: "pass1234", Hospital: "HOSP-A"})

		assert.NotEqual(t, resp1.ID, resp2.ID)
		repo.AssertExpectations(t)
	})

	t.Run("positive - same username allowed in different hospitals", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospitalB := makeHospital("HOSP-B")

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-B").Return(hospitalB, nil)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, "john_doe", hospitalB.ID.String()).Return(false, nil)
		repo.On("CreateStaff", mock.Anything, mock.MatchedBy(func(s *model.Staff) bool {
			return s.HospitalID == hospitalB.ID
		})).Return(nil)

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-B",
		})

		assert.NoError(t, err)
		assert.Equal(t, hospitalB.ID, resp.HospitalID)
		repo.AssertExpectations(t)

		// Negative unit test cast
		t.Run("negative - hospital not found returns ErrHospitalNotFound", func(t *testing.T) {
			repo := new(mockStaffRepo)
			svc := service.NewStaffService(repo)

			repo.On("FindHospitalByCode", mock.Anything, "HOSP-X").Return(nil, nil)

			resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
				Username: "john_doe", Password: "password123", Hospital: "HOSP-X",
			})

			assert.Nil(t, resp)
			assert.ErrorIs(t, err, service.ErrHospitalNotFound)
			repo.AssertNotCalled(t, "ExistsByUsernameAndHospital")
			repo.AssertNotCalled(t, "CreateStaff")
			repo.AssertExpectations(t)
		})
	})

	t.Run("negative - username already taken returns ErrUsernameAlreadyTaken", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(hospital, nil)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, "existing", hospital.ID.String()).Return(true, nil)

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "existing", Password: "password123", Hospital: "HOSP-A",
		})

		assert.Nil(t, resp)
		assert.ErrorIs(t, err, service.ErrUsernameAlreadyTaken)
		repo.AssertNotCalled(t, "CreateStaff")
		repo.AssertExpectations(t)
	})

	t.Run("negative - db error on FindHospitalByCode propagates", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(nil, errors.New("connection refused"))

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.NotErrorIs(t, err, service.ErrHospitalNotFound)
		repo.AssertNotCalled(t, "ExistsByUsernameAndHospital")
		repo.AssertExpectations(t)
	})

	t.Run("negative - db error on ExistsByUsernameAndHospital propagates", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(hospital, nil)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, "john_doe", hospital.ID.String()).
			Return(false, errors.New("db timeout"))

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		repo.AssertNotCalled(t, "CreateStaff")
		repo.AssertExpectations(t)
	})

	t.Run("negative - db error on CreateStaff propagates", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")

		repo.On("FindHospitalByCode", mock.Anything, "HOSP-A").Return(hospital, nil)
		repo.On("ExistsByUsernameAndHospital", mock.Anything, "john_doe", hospital.ID.String()).Return(false, nil)
		repo.On("CreateStaff", mock.Anything, mock.Anything).Return(errors.New("insert failed"))

		resp, err := svc.Create(context.Background(), dto.CreateStaffRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

// Test staff login
func TestStaffService_Login(t *testing.T) {
	t.Setenv("JWT_SECRET", "test_secret_at_least_32_characters_long")
	t.Setenv("JWT_EXPIRES_IN", "24h")

	// Positive unit test cast
	t.Run("positive - valid credentials returns token", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")
		staff := makeHashedStaff("john_doe", "password123", hospital.ID)

		repo.On("FindByUsernameAndHospitalCode", mock.Anything, "john_doe", "HOSP-A").Return(staff, nil)
		repo.On("UpdateLastLogin", mock.Anything, staff.ID.String(), mock.AnythingOfType("time.Time")).Return(nil)

		resp, err := svc.Login(context.Background(), dto.LoginRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.True(t, resp.ExpiresAt.After(time.Now()))
		repo.AssertExpectations(t)
	})

	t.Run("positive - login succeeds even if UpdateLastLogin fails", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")
		staff := makeHashedStaff("john_doe", "password123", hospital.ID)

		repo.On("FindByUsernameAndHospitalCode", mock.Anything, "john_doe", "HOSP-A").Return(staff, nil)
		repo.On("UpdateLastLogin", mock.Anything, staff.ID.String(), mock.AnythingOfType("time.Time")).
			Return(errors.New("db write failed"))

		resp, err := svc.Login(context.Background(), dto.LoginRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
		})

		assert.NoError(t, err)         // login ยังสำเร็จ
		assert.NotEmpty(t, resp.Token) // token ยังได้คืน
		repo.AssertExpectations(t)
	})

	// Negative unit test cast
	t.Run("negative - staff not found returns ErrInvalidCredentials", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)

		repo.On("FindByUsernameAndHospitalCode", mock.Anything, "ghost", "HOSP-A").Return(nil, nil)

		resp, err := svc.Login(context.Background(), dto.LoginRequest{
			Username: "ghost", Password: "password123", Hospital: "HOSP-A",
		})

		assert.Nil(t, resp)
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
		repo.AssertExpectations(t)
	})

	t.Run("negative - wrong password returns ErrInvalidCredentials", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)
		hospital := makeHospital("HOSP-A")
		staff := makeHashedStaff("john_doe", "correct_pass", hospital.ID)

		repo.On("FindByUsernameAndHospitalCode", mock.Anything, "john_doe", "HOSP-A").Return(staff, nil)

		resp, err := svc.Login(context.Background(), dto.LoginRequest{
			Username: "john_doe", Password: "wrong_pass", Hospital: "HOSP-A",
		})

		assert.Nil(t, resp)
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
		repo.AssertExpectations(t)
	})

	t.Run("negative - wrong hospital returns ErrInvalidCredentials", func(t *testing.T) {
		repo := new(mockStaffRepo)
		svc := service.NewStaffService(repo)

		repo.On("FindByUsernameAndHospitalCode", mock.Anything, "john_doe", "HOSP-B").Return(nil, nil)

		resp, err := svc.Login(context.Background(), dto.LoginRequest{
			Username: "john_doe", Password: "password123", Hospital: "HOSP-B",
		})

		assert.Nil(t, resp)
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
		repo.AssertExpectations(t)

		t.Run("negative - db error is NOT ErrInvalidCredentials", func(t *testing.T) {
			repo := new(mockStaffRepo)
			svc := service.NewStaffService(repo)

			repo.On("FindByUsernameAndHospitalCode", mock.Anything, "john_doe", "HOSP-A").
				Return(nil, errors.New("connection refused"))

			resp, err := svc.Login(context.Background(), dto.LoginRequest{
				Username: "john_doe", Password: "password123", Hospital: "HOSP-A",
			})

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.NotErrorIs(t, err, service.ErrInvalidCredentials)
			repo.AssertExpectations(t)
		})
	})
}
