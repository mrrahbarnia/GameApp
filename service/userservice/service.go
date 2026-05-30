package userservice

import (
	"fmt"

	"github.com/mrrahbarnia/GameApp/entity"
	"github.com/mrrahbarnia/GameApp/pkg/phonenumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
}

type Service struct {
	repo Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
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
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("Unexpected error: %w", err)
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("Phone number is not unique")
		}
	}

	// Validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("Name length must be greater than 3")
	}

	// Create user
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
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
