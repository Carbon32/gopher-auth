package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/mail"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUser
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("A user with the e-mail %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: "",
		LastName:  "",
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  hashedPassword,
		Verified:  false,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

	err = mail.SendEmail(payload.Email, "Welcome", mail.VerificationTemplate(payload.Username, "www.whatever.com"), true)

	if err != nil {
		log.Println("EMAIL ERROR: ", err)
	}

}
