package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value(auth.UserIdKey).(int)

	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("User not found"))
		return
	}

	response := types.Profile{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Verified:  user.Verified,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(auth.UserIdKey).(int)

	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	var payload types.UpdateProfile
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !regexp.MustCompile(`^[\p{L} ]+$`).MatchString(payload.FirstName) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("First name may only contain letters and spaces"))
		return
	}

	if !regexp.MustCompile(`^[\p{L} ]+$`).MatchString(payload.LastName) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Last name may only contain letters and spaces"))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(payload.Username) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Username may only contain letters, numbers, and underscores"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request %v", errors))
		return
	}

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Username = payload.Username

	err = h.store.UpdateUser(*user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Information updated")
}

func (h *Handler) updateEmail(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(auth.UserIdKey).(int)

	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	var payload types.UpdateEmail
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request %v", errors))
		return
	}

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	if !auth.ComparePassword(user.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid password"))
		return
	}

	if user.Email == payload.Email {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("The user already uses this e-mail"))
		return
	}

	user.Email = payload.Email

	err = h.store.UpdateUser(*user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Information updated")
}

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(auth.UserIdKey).(int)

	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	var payload types.UpdatePassword
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request %v", errors))
		return
	}

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	if !auth.ComparePassword(user.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid password"))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !auth.ComparePassword(user.Password, []byte(hashedPassword)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("The user already uses this password"))
		return
	}

	user.Password = hashedPassword

	err = h.store.UpdateUser(*user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Information updated")
}
