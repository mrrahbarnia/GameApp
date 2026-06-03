package userservice

import (
	"fmt"

	"github.com/mrrahbarnia/GameApp/entity"
	"github.com/mrrahbarnia/GameApp/pkg/richerror"
	"github.com/mrrahbarnia/GameApp/presentation/dto"
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

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	if exist, err := s.repo.IsPhoneNumberExist(req.PhoneNumber); err != nil || exist {
		if err != nil {
			return dto.RegisterResponse{},
				richerror.New("userservice.Register").
					WithErr(err).
					WithKind(richerror.KindUnexpected)
		}

		if exist {
			return dto.RegisterResponse{},
				richerror.New("userservice.Register").
					WithKind(richerror.KindConflict).
					WithMessage("phone_number is already exist").
					WithMeta(map[string]any{"phone_number": req.PhoneNumber})
		}
	}

	hashedPassword, err := s.bcrypt.GeneratePasswordHash(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("Unexpected error: %w", err)
	}

	user := entity.User{
		ID:             0,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		HashedPassword: hashedPassword,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("Unexpected error: %w", err)
	}

	return dto.RegisterResponse{
		UserID:      createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}, nil

}

// ******************************** Login usecase

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	dbUser, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("Unexpected error: %w", err)
	}
	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("Wrong credentials")
	}

	if !s.bcrypt.ComparePassword(dbUser.HashedPassword, req.Password) {
		return dto.LoginResponse{}, fmt.Errorf("Wrong credentials")
	}

	accessToken, aErr := s.auth.CreateAccessToken(dbUser)
	refreshToken, rErr := s.auth.CreateRefreshToken(dbUser)
	if aErr != nil {
		return dto.LoginResponse{}, fmt.Errorf("Unexpected error: %w", aErr)
	}
	if rErr != nil {
		return dto.LoginResponse{}, fmt.Errorf("Unexpected error: %w", rErr)
	}

	return dto.LoginResponse{
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
