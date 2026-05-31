package userservice

import (
	"fmt"

	"github.com/mrrahbarnia/GameApp/entity"
	"github.com/mrrahbarnia/GameApp/pkg/phonenumber"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberExist(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func comparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// ******************************** Register usecase

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - We should verify phone number by verification code
	// Check phone number validity
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("Phone number is not valid")
	}

	// Check phone number uniqeness
	if exist, err := s.repo.IsPhoneNumberExist(req.PhoneNumber); err != nil || exist {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("Unexpected error: %w", err)
		}

		if exist {
			return RegisterResponse{}, fmt.Errorf("Phone number is not unique")
		}
	}

	// Validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("Name length must be greater than 3")
	}

	// Validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("Password length must be greater than 8")
	}
	hashedPassword, err := generatePasswordHash(req.Password)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("Unexpected error: %w", err)
	}

	// Create user
	user := entity.User{
		ID:             0,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		HashedPassword: hashedPassword,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("Unexpected error: %w", err)
	}

	// Response
	return RegisterResponse{
		User: createdUser,
	}, nil

}

// ******************************** Login usecase

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	User entity.User
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// Check phone number validity
	if !phonenumber.IsValid(req.PhoneNumber) {
		return LoginResponse{}, fmt.Errorf("Phone number is not valid")
	}

	// Check user exist
	dbUser, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Unexpected error: %w", err)
	}
	if !exist {
		return LoginResponse{}, fmt.Errorf("Wrong credentials")
	}

	if !comparePassword(dbUser.HashedPassword, req.Password) {
		return LoginResponse{}, fmt.Errorf("Wrong credentials")
	}

	return LoginResponse{}, nil
}
