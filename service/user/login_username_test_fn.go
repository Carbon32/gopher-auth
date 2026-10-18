package user

import (
	"bytes"
	"carb-backend/service/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func packAndSendLoginUsernamePayload(payload types.LoginUserUsername, t *testing.T, handler *Handler, expected int) {
	json_payload, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/login_username", bytes.NewBuffer(json_payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/login_username", handler.handleLoginUsername)
	router.ServeHTTP(rr, req)

	if rr.Code != expected {
		t.Errorf("Expected status code %d, got %d: %s", expected, rr.Code, rr.Body.String())
	}
}

func testLoginUsernameValidRequest(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should pass if the request is valid", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   username,
			Password:   password,
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusOK)
	})
}

func testLoginUsernameInvalidRequest(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should not pass if the request is contains an incorrect password", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   username,
			Password:   "123123123123",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusUnauthorized)
	})
}

func testLoginUsernameMissingFields(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the request is missing a field", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username: username,
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginUsernameAllFieldsEmpty(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the request fields are empty", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   "",
			Password:   "",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginUsernameSomeFieldsAreEmpty(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the request fields are empty", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   username,
			Password:   "",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginUsernameUsernameTooShort(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the username is too short", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   "ab",
			Password:   "whatever",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginUsernameUsernameTooLong(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the username is too long", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   "pneumonoultramicroscopicsilcicovolcanoconiosis",
			Password:   "whatever",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginUsernameBadUsername(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the username contains unwanted characters", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   "carbon@",
			Password:   "whatever",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginUsernamePasswordTooShort(t *testing.T, handler *Handler) {
	t.Run("LOGIN (USERNAME): This should fail if the password is short", func(t *testing.T) {
		payload := types.LoginUserUsername{
			Username:   username,
			Password:   "abc",
			RememberMe: new(false),
		}
		packAndSendLoginUsernamePayload(payload, t, handler, http.StatusBadRequest)
	})
}
