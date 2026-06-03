package userservice

import (
	"fmt"

	"github.com/mrrahbarnia/GameApp/entity"
	"github.com/mrrahbarnia/GameApp/pkg/phonenumber"
	"github.com/mrrahbarnia/GameApp/pkg/richerror"
)

type Repository interface {
	IsPhoneNumberExist(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserById(userId uint) (entity.User, bool, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Bcrypt interface {
	GeneratePasswordHash(password string) (string, error)
	ComparePassword(hashedPassword, plainPassword string) bool
}

type Service struct {
	auth   AuthGenerator
	repo   Repository
	bcrypt Bcrypt
}

func New(repo Repository, bcrypt Bcrypt, authGenerator AuthGenerator) Service {
	return Service{repo: repo, bcrypt: bcrypt, auth: authGenerator}
}

// ******************************** Register usecase

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - We should verify phone number by verification code
	// Check phone number validity
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{},
			richerror.New("userservice.Register").
				WithKind(richerror.KindInvalid).
				WithMessage("phone_number is not valid").
				WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	// Check phone number uniqeness
	if exist, err := s.repo.IsPhoneNumberExist(req.PhoneNumber); err != nil || exist {
		if err != nil {
			return RegisterResponse{},
				richerror.New("userservice.Register").
					WithErr(err).
					WithKind(richerror.KindUnexpected)
		}

		if exist {
			return RegisterResponse{},
				richerror.New("userservice.Register").
					WithKind(richerror.KindConflict).
					WithMessage("phone_number is already exist").
					WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
		}
	}

	// Validate name
	if len(req.Name) < 3 {
		return RegisterResponse{},
			fmt.Errorf("Name length must be greater than 3")
	}

	// Validate password
	if len(req.Password) < 8 {
		return RegisterResponse{},
			fmt.Errorf("Password length must be greater than 8")
	}
	hashedPassword, err := s.bcrypt.GeneratePasswordHash(req.Password)
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
		UserID:      createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}, nil

}

// ******************************** Login usecase

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	if !s.bcrypt.ComparePassword(dbUser.HashedPassword, req.Password) {
		return LoginResponse{}, fmt.Errorf("Wrong credentials")
	}

	accessToken, aErr := s.auth.CreateAccessToken(dbUser)
	refreshToken, rErr := s.auth.CreateRefreshToken(dbUser)
	if aErr != nil {
		return LoginResponse{}, fmt.Errorf("Unexpected error: %w", aErr)
	}
	if rErr != nil {
		return LoginResponse{}, fmt.Errorf("Unexpected error: %w", rErr)
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ******************************** Get profile usecase

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	dbUser, exist, err := s.repo.GetUserById(req.UserID)
	if err != nil {
		return ProfileResponse{},
			richerror.New("userservice.Profile").WithErr(err).WithMeta(map[string]interface{}{"req": req})
	}
	if !exist {
		return ProfileResponse{},
			richerror.New("userservice.Profile").WithKind(richerror.KindNotFound)
	}

	return ProfileResponse{
		Name: dbUser.Name,
	}, nil

}
