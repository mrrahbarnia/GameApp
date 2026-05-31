package bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
}

func New() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func (b *Bcrypt) ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
