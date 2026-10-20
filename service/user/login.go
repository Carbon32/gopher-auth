package user

import (
	"carb-backend/cmd/db"
	"carb-backend/service/auth"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var dummyPasswordHash, _ = auth.HashPassword("not-an-account-password")

func (h *Handler) login(w http.ResponseWriter, r *http.Request, lookup func(string) (*types.User, error), identifier, password string, remember bool) {
	u, err := lookup(identifier)
	if errors.Is(err, ErrNotFound) {
		_ = auth.ComparePassword(dummyPasswordHash, []byte(password))
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Invalid credentials"))
		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, fmt.Errorf("Authentication service unavailable"))
		return
	}
	if !auth.ComparePassword(u.Password, []byte(password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Invalid credentials"))
		return
	}
	if !u.Verified {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("Verify your e-mail before signing in"))
		return
	}
	sid, err := auth.RandomID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	refresh, err := auth.RandomToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	lifetime := 24 * time.Hour
	if remember {
		lifetime = 30 * 24 * time.Hour
	}
	sessionExpiresAt := time.Now().Add(lifetime)
	if err := h.store.CreateSession(r.Context(), sid, u.Id, hashToken(refresh), time.Now().Add(lifetime)); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not create session"))
		return
	}
	access, err := auth.CreateAccessJWT([]byte(db.Environement.JWTSecret), u.Id, sid)
	if err != nil {
		_ = h.store.RevokeSession(context.Background(), sid)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": refresh, "expires_in": 900, "expires_at": sessionExpiresAt.UTC()})
}

func (h *Handler) handleLoginEmail(w http.ResponseWriter, r *http.Request) {
	var p types.LoginUserEmail
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}
	h.login(w, r, h.store.GetUserByEmail, strings.ToLower(strings.TrimSpace(p.Email)), p.Password, *p.RememberMe)
}

func (h *Handler) handleLoginUsername(w http.ResponseWriter, r *http.Request) {
	var p types.LoginUserUsername
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}
	if !usernamePattern.MatchString(p.Username) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}
	h.login(w, r, h.store.GetUserByUsername, strings.ToLower(strings.TrimSpace(p.Username)), p.Password, *p.RememberMe)
}
