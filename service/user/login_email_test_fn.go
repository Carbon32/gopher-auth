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

func packAndSendLoginEmailPayload(payload types.LoginUserEmail, t *testing.T, handler *Handler, expected int) {
	json_payload, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/login_email", bytes.NewBuffer(json_payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/login_email", handler.handleLoginEmail)
	router.ServeHTTP(rr, req)

	if rr.Code != expected {
		t.Errorf("Expected status code %d, got %d", expected, rr.Code)
	}
}

func testLoginEmailValidRequest(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should pass if the request is valid", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email:      email,
			Password:   password,
			RememberMe: new(false),
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusOK)
	})
}

func testLoginEmailInvalidRequest(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should not pass if the request is contains an incorrect password", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email:      email,
			Password:   "123123123123",
			RememberMe: new(false),
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusUnauthorized)
	})
}

func testLoginEmailMissingFields(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should fail if the request is missing a field", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email: email,
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginEmailAllFieldsEmpty(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should fail if the request fields are empty", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email:      "",
			Password:   "",
			RememberMe: new(false),
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginEmailSomeFieldsAreEmpty(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should fail if the request fields are empty", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email:      email,
			Password:   "",
			RememberMe: new(false),
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginEmailBadEmail(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should fail if the e-mail is invlaid", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email:      "carbon@some_bad_stuff",
			Password:   "whatever",
			RememberMe: new(false),
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusBadRequest)
	})
}

func testLoginEmailPasswordTooShort(t *testing.T, handler *Handler) {
	t.Run("LOGIN (E-MAIL): This should fail if the password is short", func(t *testing.T) {
		payload := types.LoginUserEmail{
			Email:      email,
			Password:   "abc",
			RememberMe: new(false),
		}
		packAndSendLoginEmailPayload(payload, t, handler, http.StatusBadRequest)
	})
}
