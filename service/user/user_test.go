package user

import (
	"testing"
)

func TestUserServiceHandlers(t *testing.T) {
	registerUserStore := &registerTestUserStore{}
	handler := NewHandler(registerUserStore)

	// Register
	testValidRequest(t, handler)
	testMissingFields(t, handler)
	testAllFieldsEmpty(t, handler)
	testSomeFieldsAreEmpty(t, handler)
	testUsernameTooShort(t, handler)
	testUsernameTooLong(t, handler)
	testBadUsername(t, handler)
	testInvalidEmail(t, handler)
	testPasswordTooShort(t, handler)

	loginUserStore := &loginTestUserStore{}
	handler = NewHandler(loginUserStore)

	// Login (username)
	testLoginUsernameValidRequest(t, handler)
	testLoginUsernameInvalidRequest(t, handler)
	testLoginUsernameMissingFields(t, handler)
	testLoginUsernameAllFieldsEmpty(t, handler)
	testLoginUsernameSomeFieldsAreEmpty(t, handler)
	testLoginUsernameUsernameTooShort(t, handler)
	testLoginUsernameUsernameTooLong(t, handler)
	testLoginUsernameBadUsername(t, handler)
	testLoginUsernamePasswordTooShort(t, handler)

	// Login (e-mail)
	testLoginEmailValidRequest(t, handler)
	testLoginEmailInvalidRequest(t, handler)
	testLoginEmailMissingFields(t, handler)
	testLoginEmailAllFieldsEmpty(t, handler)
	testLoginEmailSomeFieldsAreEmpty(t, handler)
	testLoginEmailBadEmail(t, handler)
	testLoginEmailPasswordTooShort(t, handler)
}
