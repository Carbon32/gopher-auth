package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/types"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type registerTestUserStore struct{}
type loginTestUserStore struct{}

// This account should exist for all login test cases to fire correctly, you can change it to whatever you want

const username = "carbon32"
const email = "carbon@test.com"
const password = "hello_world123"

func TestUserStruct() *types.User {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		panic(err)
	}

	testUser := types.User{
		Id:        1,
		FirstName: "",
		LastName:  "",
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		Verified:  true,
	}

	return &testUser
}

func (t *registerTestUserStore) CreateUser(types.User) error {
	return nil
}

func (t *registerTestUserStore) UpdateUser(types.User) error {
	return nil
}

func (t *registerTestUserStore) GetUserById(id int) (*types.User, error) {
	return nil, fmt.Errorf("User not found")
}

func (t *registerTestUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, sql.ErrNoRows
}

func (t *registerTestUserStore) GetUserByUsername(username string) (*types.User, error) {
	return nil, sql.ErrNoRows
}

func (t *registerTestUserStore) ChangePassword(context.Context, int, string) error { return nil }
func (t *registerTestUserStore) CreateActionToken(context.Context, int, string, string, string, time.Time) error {
	return nil
}
func (t *registerTestUserStore) CompleteEmailVerification(context.Context, string) error { return nil }
func (t *registerTestUserStore) CompleteEmailChange(context.Context, string) error       { return nil }
func (t *registerTestUserStore) CompletePasswordReset(context.Context, string, string) error {
	return nil
}
func (t *registerTestUserStore) CreateSession(context.Context, string, int, string, time.Time) error {
	return nil
}
func (t *registerTestUserStore) SessionActive(context.Context, string) (bool, error) {
	return true, nil
}
func (t *registerTestUserStore) RotateRefreshToken(context.Context, string, string) (string, int, time.Time, error) {
	return "", 0, time.Time{}, sql.ErrNoRows
}
func (t *registerTestUserStore) RevokeSession(context.Context, string) error { return nil }

func (t *loginTestUserStore) ChangePassword(context.Context, int, string) error           { return nil }
func (t *loginTestUserStore) CompleteEmailVerification(context.Context, string) error     { return nil }
func (t *loginTestUserStore) CompleteEmailChange(context.Context, string) error           { return nil }
func (t *loginTestUserStore) CompletePasswordReset(context.Context, string, string) error { return nil }
func (t *loginTestUserStore) CreateSession(context.Context, string, int, string, time.Time) error {
	return nil
}
func (t *loginTestUserStore) SessionActive(context.Context, string) (bool, error) { return true, nil }
func (t *loginTestUserStore) RotateRefreshToken(context.Context, string, string) (string, int, time.Time, error) {
	return "", 0, time.Time{}, sql.ErrNoRows
}
func (t *loginTestUserStore) RevokeSession(context.Context, string) error { return nil }

func (t *loginTestUserStore) CreateUser(types.User) error {
	return nil
}

func (t *loginTestUserStore) UpdateUser(types.User) error {
	return nil
}

func (t *loginTestUserStore) GetUserById(id int) (*types.User, error) {
	return TestUserStruct(), nil
}

func (t *loginTestUserStore) GetUserByEmail(email string) (*types.User, error) {
	return TestUserStruct(), nil
}

func (t *loginTestUserStore) GetUserByUsername(username string) (*types.User, error) {
	return TestUserStruct(), nil
}

func (t *loginTestUserStore) CreateActionToken(context.Context, int, string, string, string, time.Time) error {
	return nil
}
