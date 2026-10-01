package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/types"
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
		Verified:  false,
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
	return nil, fmt.Errorf("User not found")
}

func (t *registerTestUserStore) GetUserByUsername(username string) (*types.User, error) {
	return nil, fmt.Errorf("User not found")
}

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
