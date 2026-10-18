package user

import (
	"carb-backend/cmd/db"
	"carb-backend/service/auth"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	limiter := auth.NewLoginLimiter(10, time.Minute)
	router.Handle("/login_email", limiter.Middleware(http.HandlerFunc(h.handleLoginEmail))).Methods("POST")
	router.Handle("/login_username", limiter.Middleware(http.HandlerFunc(h.handleLoginUsername))).Methods("POST")
	router.HandleFunc("/verify_email", h.verifyEmail).Methods("GET", "POST")
	router.HandleFunc("/email/change/confirm", h.confirmEmailChange).Methods("GET", "POST")
	router.Handle("/email/verification/resend", auth.NewLoginLimiter(3, time.Minute).Middleware(http.HandlerFunc(h.resendVerification))).Methods("POST")
	router.Handle("/password/forgot", auth.NewLoginLimiter(3, time.Minute).Middleware(http.HandlerFunc(h.forgotPassword))).Methods("POST")
	router.HandleFunc("/password/reset", h.resetPassword).Methods("POST")
	router.HandleFunc("/session/refresh", h.refreshSession).Methods("POST")

	protected := router.PathPrefix("/user").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.RequireAuth([]byte(db.Environement.JWTSecret), h.store.SessionActive, next)
	})

	protected.HandleFunc("/me", h.handleMe).Methods("GET")
	protected.HandleFunc("/logout", h.logout).Methods("POST")

	protected.HandleFunc("/profile", h.getProfile).Methods("GET")
	protected.HandleFunc("/update_profile", h.updateProfile).Methods("PUT")
	protected.HandleFunc("/update_email", h.updateEmail).Methods("PUT")
	protected.HandleFunc("/update_password", h.updatePassword).Methods("PUT")

}
