package user

import (
	"carb-backend/cmd/db"
	"carb-backend/service/auth"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login_email", h.handleLoginEmail).Methods("POST")
	router.HandleFunc("/login_username", h.handleLoginUsername).Methods("POST")

	protected := router.PathPrefix("/user").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.RequireAuth(
			[]byte(db.Environement.JWTSecret),
			next,
		)
	})

	protected.HandleFunc("/me", h.handleMe).Methods("GET")
	protected.HandleFunc("/confirm", h.handleConfirm).Methods("GET")

	protected.HandleFunc("/profile", h.getProfile).Methods("GET")
	protected.HandleFunc("/update_profile", h.updateProfile).Methods("PUT")
	protected.HandleFunc("/update_email", h.updateEmail).Methods("PUT")
	protected.HandleFunc("/update_passowrd", h.updatePassword).Methods("PUT")

}
