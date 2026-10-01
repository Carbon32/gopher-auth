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

func packAndSendRegisterPayload(payload types.RegisterUser, t *testing.T, handler *Handler, expected int) {
	json_payload, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(json_payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/register", handler.handleRegister)
	router.ServeHTTP(rr, req)

	if rr.Code != expected {
		t.Errorf("Expected status code %d, got %d: %s", expected, rr.Code, rr.Body.String())
	}
}

func testValidRequest(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should pass if the request is valid", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "carbon",
			Email:    "new_email@domain.com",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusCreated)
	})
}

func testMissingFields(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the request is missing a field", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "carbon",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testAllFieldsEmpty(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the request fields are empty", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "",
			Email:    "",
			Password: "",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testSomeFieldsAreEmpty(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the request fields are empty", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "",
			Email:    "new_email@domain.com",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testUsernameTooShort(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the username is too short", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "ab",
			Email:    "new_email@domain.com",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testUsernameTooLong(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the username is too long", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "pneumonoultramicroscopicsilcicovolcanoconiosis",
			Email:    "new_email@domain.com",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testBadUsername(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the username contains unwanted characters", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "carbon@",
			Email:    "new_email@domain.com",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testInvalidEmail(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the e-mail is invlaid", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "carbon",
			Email:    "new_emaildomain.com",
			Password: "whatever",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testPasswordTooShort(t *testing.T, handler *Handler) {
	t.Run("REGISTER: This should fail if the password is short", func(t *testing.T) {
		payload := types.RegisterUser{
			Username: "carbon",
			Email:    "new_email@domain.com",
			Password: "abc",
		}
		packAndSendRegisterPayload(payload, t, handler, http.StatusBadRequest)
	})
}
