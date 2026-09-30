package types

import "time"

type UserStore interface {
	CreateUser(User) error
	UpdateUser(User) error
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

type User struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Verified  bool      `json:"verified"`
}

type RegisterUser struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type LoginUserEmail struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe *bool  `json:"remember_me" validate:"required"`
}

type LoginUserUsername struct {
	Username   string `json:"username" validate:"required,min=3,max=32"`
	Password   string `json:"password" validate:"required"`
	RememberMe *bool  `json:"remember_me" validate:"required"`
}

type Profile struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Verified  bool      `json:"verified"`
}

type UpdateProfile struct {
	FirstName string `json:"firstName" validate:"required,min=3,max=32"`
	LastName  string `json:"lastName" validate:"required,min=3,max=32"`
	Username  string `json:"username" validate:"required,min=3,max=32"`
}

type UpdateEmail struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type UpdatePassword struct {
	Password    string `json:"password" validate:"required,min=8,max=128"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=128"`
}

// Whatever, I'm noting things here
// profile (x) t ()
// update_profile (x) t ()
// update_email (x) t ()
// update_password (x) t ()
// mailpit t ()
