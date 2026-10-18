package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuthChecksSignatureAndSession(t *testing.T) {
	secret := []byte("test-secret")
	token, err := CreateAccessJWT(secret, 12, "session-1")
	if err != nil {
		t.Fatal(err)
	}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if got := r.Context().Value(UserIdKey); got != 12 {
			t.Errorf("user id = %v, want 12", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	checker := func(_ context.Context, id string) (bool, error) { return id == "session-1", nil }
	handler := RequireAuth(secret, checker, next)
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNoContent || !called {
		t.Fatalf("response = %d, next called = %v", recorder.Code, called)
	}

	called = false
	revoked := RequireAuth(secret, func(context.Context, string) (bool, error) { return false, nil }, next)
	recorder = httptest.NewRecorder()
	revoked.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusUnauthorized || called {
		t.Fatalf("revoked token response = %d, next called = %v", recorder.Code, called)
	}
}
