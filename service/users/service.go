package userservice

import (
	"fmt"

	"github.com/mrrahbarnia/GameApp/entity"
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

type RegisterOut struct {
	UserID      uint
	Name        string
	PhoneNumber string
}

func (s Service) Register(name, phoneNumber, password string) (RegisterOut, error) {
	if exist, err := s.repo.IsPhoneNumberExist(phoneNumber); err != nil || exist {
		if err != nil {

			return RegisterOut{},
				richerror.New("userservice.Register").
					WithErr(err).
					WithKind(richerror.KindUnexpected)
		}

		if exist {
			return RegisterOut{},
				richerror.New("userservice.Register").
					WithKind(richerror.KindConflict).
					WithMessage("phone_number is already exist").
					WithMeta(map[string]any{"phone_number": phoneNumber})
		}
	}

	hashedPassword, err := s.bcrypt.GeneratePasswordHash(password)
	if err != nil {
		return RegisterOut{}, fmt.Errorf("Unexpected error: %w", err)
	}

	user := entity.User{
		ID:             0,
		Name:           name,
		PhoneNumber:    phoneNumber,
		HashedPassword: hashedPassword,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterOut{}, fmt.Errorf("Unexpected error: %w", err)
	}

	return RegisterOut{
		UserID:      createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}, nil

}

// ******************************** Login usecase

type LoginOut struct {
	AccessToken  string
	RefreshToken string
}

func (s Service) Login(phoneNumber, password string) (LoginOut, error) {
	dbUser, exist, err := s.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return LoginOut{}, fmt.Errorf("Unexpected error: %w", err)
	}
	if !exist {
		return LoginOut{}, fmt.Errorf("Wrong credentials")
	}

	if !s.bcrypt.ComparePassword(dbUser.HashedPassword, password) {
		return LoginOut{}, fmt.Errorf("Wrong credentials")
	}

	accessToken, aErr := s.auth.CreateAccessToken(dbUser)
	refreshToken, rErr := s.auth.CreateRefreshToken(dbUser)
	if aErr != nil {
		return LoginOut{}, fmt.Errorf("Unexpected error: %w", aErr)
	}
	if rErr != nil {
		return LoginOut{}, fmt.Errorf("Unexpected error: %w", rErr)
	}

	return LoginOut{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ******************************** Get profile usecase

type ProfileOut struct {
	Name string
}

func (s Service) Profile(userId uint) (ProfileOut, error) {
	dbUser, exist, err := s.repo.GetUserById(userId)
	if err != nil {
		return ProfileOut{},
			richerror.New("userservice.Profile").WithErr(err).WithMeta(map[string]interface{}{"req": userId})
	}
	if !exist {
		return ProfileOut{},
			richerror.New("userservice.Profile").WithKind(richerror.KindNotFound)
	}

	return ProfileOut{
		Name: dbUser.Name,
	}, nil

}
