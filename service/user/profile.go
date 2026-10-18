package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	user, err := h.store.GetUserById(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}

	payload.Username = strings.ToLower(strings.TrimSpace(payload.Username))
	user.FirstName = strings.TrimSpace(payload.FirstName)
	user.LastName = payload.LastName
	user.Username = payload.Username

	err = h.store.UpdateUser(*user)
	if err != nil {
		if isUniqueViolation(err) {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("Username is unavailable"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not update profile"))
		}
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
	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
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

	if strings.EqualFold(user.Email, payload.Email) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("The account already uses this e-mail"))
		return
	}
	if _, err = h.store.GetUserByEmail(payload.Email); err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("E-mail is unavailable"))
		return
	} else if !errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not check e-mail availability"))
		return
	}
	if err := h.issueAction(r.Context(), user.Id, payload.Email, "email_change", "/email/change/confirm", "Confirm your new e-mail"); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not send confirmation e-mail"))
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, map[string]string{"message": "Confirm the new e-mail address using the link we sent"})
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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
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

	if auth.ComparePassword(user.Password, []byte(payload.NewPassword)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("New password must differ from the current password"))
		return
	}
	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.ChangePassword(r.Context(), userId, hashedPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not update password"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password updated; please log in again"})
}
