package user

import (
	"carb-backend/cmd/db"
	"carb-backend/service/auth"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func (h *Handler) handleLoginEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserEmail
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid e-mail or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid e-mail or password"))
		return
	}

	token, err := auth.CreateJWT([]byte(db.Environement.JWTSecret), u.Id, *payload.RememberMe)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": token})
}

func (h *Handler) handleLoginUsername(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserUsername
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request %v", errors))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(payload.Username) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Username may only contain letters, numbers, and underscores"))
		return
	}

	u, err := h.store.GetUserByUsername(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid username or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid username or password"))
		return
	}

	token, err := auth.CreateJWT([]byte(db.Environement.JWTSecret), u.Id, *payload.RememberMe)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"access_token": token})
}
