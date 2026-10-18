package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var p types.RegisterUser
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	if !usernamePattern.MatchString(p.Username) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Username may only contain letters, numbers, and underscores"))
		return
	}
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}
	if len([]byte(p.Password)) > 72 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Password must be at most 72 bytes"))
		return
	}
	if _, err := h.store.GetUserByEmail(p.Email); err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("E-mail or username is unavailable"))
		return
	} else if !errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not check account availability"))
		return
	}
	if _, err := h.store.GetUserByUsername(p.Username); err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("E-mail or username is unavailable"))
		return
	} else if !errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not check account availability"))
		return
	}
	hash, err := auth.HashPassword(p.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.CreateUser(types.User{Username: p.Username, Email: p.Email, Password: hash, Verified: false}); err != nil {
		if isUniqueViolation(err) {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("E-mail or username is unavailable"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not create account"))
		}
		return
	}
	u, err := h.store.GetUserByEmail(p.Email)
	if err == nil {
		if e := h.issueAction(r.Context(), u.Id, p.Email, "verify_email", "/verify_email", "Verify your e-mail"); e != nil {
			log.Printf("verification e-mail error: %v", e)
		}
	} else {
		log.Printf("could not load newly created user for verification: %v", err)
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Account created. Check your e-mail to verify it."})
}
