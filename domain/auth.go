package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
	"time"
)

type AuthStore interface {
	Create(u SignUpInput) (*User, error)
	GetOne(email string, id string) (*User, error)
	SetAccessToken(uuid uuid.UUID) (*string, error)
	SetRefreshToken(uuid uuid.UUID) (*string, error)
	GetByRefreshTokenRedis(token string) (string, error)
	GetByAccessTokenRedis(token string) (string, string, error)
	DeleteTokensRedis(refreshToken, accessToken string) error
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SignUpInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required,min=8"`
}

func (r *SignUpInput) HashPassword() error {
	if len(r.Password) == 0 {
		return errors.New("password should not be empty")
	}

	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(r.Password))

	if err != nil {
		return err
	}
	r.Password = string(encoded)

	return nil
}

func (r *SignInInput) CheckPassword(plain string) (bool, error) {
	return argon2.VerifyEncoded([]byte(r.Password), []byte(plain))
}
