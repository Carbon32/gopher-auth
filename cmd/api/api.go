package api

// I'm gonna keep this an authentication system only, I'm gonna build the rest of the system using C# or Python, not sure

import (
	"carb-backend/service/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	address *string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: &address,
		db:      db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // For testing
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}).Handler(router)

	log.Println("Listening on", *s.address)

	return http.ListenAndServe(*s.address, handler)
}
