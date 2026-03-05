package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"twichai/agnos-test/api/middleware"
	"twichai/agnos-test/api/presenter/patient"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/repository"
	"twichai/agnos-test/pkg/staff/usecase"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffService struct {
	repo      repository.StaffRepository
	jwtSecret []byte
}

func NewStaffService(repo repository.StaffRepository, appSecret string) usecase.StaffUseCase {
	return &StaffService{repo: repo, jwtSecret: []byte(appSecret)}
}

// Create implements [usecase.StaffUseCase].
func (s *StaffService) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	hashedPassword, err := hashPassword(staff.Password)
	if err != nil {
		return nil, err
	}

	staff.Password = hashedPassword
	return s.repo.Create(ctx, staff)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Login implements [usecase.StaffUseCase].
func (s *StaffService) Login(ctx context.Context, staff *entity.StaffLoginRequest) (*patient.LoginStaffPresenter, error) {
	existingStaff, err := s.repo.GetByUsername(ctx, staff.Username, staff.HospitalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usecase.ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingStaff.Password), []byte(staff.Password))
	if err != nil {
		return nil, usecase.ErrInvalidCredentials
	}

	if len(s.jwtSecret) == 0 {
		return nil, errors.New("jwt app secret is not configured")
	}

	generatedToken, err := s.generateToken(existingStaff)
	if err != nil {
		fmt.Printf("failed to generate token: %v\n", err)
		return nil, err
	}

	return &patient.LoginStaffPresenter{
		Token: generatedToken,
	}, nil

}

func (s *StaffService) generateToken(staff *entity.Staff) (string, error) {
	now := time.Now()

	claims := middleware.Claims{
		StaffName:  staff.Username,
		HospitalID: staff.Hospital.ID,
		StaffID:    staff.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   staff.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
