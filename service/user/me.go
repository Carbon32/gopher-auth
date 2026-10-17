package user

import (
	"carb-backend/service/auth"
	"carb-backend/service/utils"
	"fmt"
	"net/http"
)

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIdKey).(int)

	if !ok {
		utils.WriteError(
			w,
			http.StatusUnauthorized,
			fmt.Errorf("Unauthorized"),
		)
		return
	}

	u, err := h.store.GetUserById(userID)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusNotFound,
			fmt.Errorf("User not found"),
		)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"id":         u.Id,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"username":   u.Username,
		"email":      u.Email,
	})
}

func (h *Handler) handleConfirm(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIdKey).(int)

	if !ok {
		utils.WriteError(
			w,
			http.StatusUnauthorized,
			fmt.Errorf("Unauthorized"),
		)
		return
	}

	u, err := h.store.GetUserById(userID)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusNotFound,
			fmt.Errorf("User not found"),
		)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Welcome " + u.Username})
}
