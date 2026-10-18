package types

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserStore interface {
	CreateUser(User) error
	UpdateUser(User) error
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	ChangePassword(ctx context.Context, id int, passwordHash string) error
	CreateActionToken(ctx context.Context, userID int, purpose, targetEmail, tokenHash string, expiresAt time.Time) error
	CompleteEmailVerification(ctx context.Context, tokenHash string) error
	CompleteEmailChange(ctx context.Context, tokenHash string) error
	CompletePasswordReset(ctx context.Context, tokenHash, passwordHash string) error
	CreateSession(ctx context.Context, id string, userID int, refreshHash string, expiresAt time.Time) error
	SessionActive(ctx context.Context, id string) (bool, error)
	RotateRefreshToken(ctx context.Context, oldRefreshHash, newRefreshHash string) (string, int, time.Time, error)
	RevokeSession(ctx context.Context, id string) error
}

type User struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Verified  bool      `json:"verified"`
}

type AccessClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type SessionChecker func(context.Context, string) (bool, error)

type RegisterUser struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginUserEmail struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
	RememberMe *bool  `json:"remember_me" validate:"required"`
}

type LoginUserUsername struct {
	Username   string `json:"username" validate:"required,min=3,max=32"`
	Password   string `json:"password" validate:"required,min=8,max=72"`
	RememberMe *bool  `json:"remember_me" validate:"required"`
}

type Profile struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Verified  bool      `json:"verified"`
}

type UpdateProfile struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string `json:"last_name" validate:"required,min=3,max=32"`
	Username  string `json:"username" validate:"required,min=3,max=32"`
}

type UpdateEmail struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}
type ResetPassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

type UpdatePassword struct {
	Password    string `json:"password" validate:"required,min=8,max=72"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}
