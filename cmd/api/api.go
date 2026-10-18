package api

import (
	"carb-backend/cmd/db"
	"carb-backend/service/user"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, database *sql.DB) *APIServer {
	if address == "" {
		address = db.Environement.HTTPAddress
	}
	return &APIServer{address: address, db: database}
}

// A little bit of security won't hurt
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) handler() http.Handler {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	user.NewHandler(user.NewStore(s.db)).RegisterRoutes(subrouter)
	return securityHeaders(cors.New(cors.Options{
		AllowedOrigins:   db.Environement.CORSAllowedOrigins, // Cleaner than just throwing an URL here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router))
}

func (s *APIServer) Run() error {
	return s.RunContext(context.Background())
}

func (s *APIServer) RunContext(ctx context.Context) error {
	if (db.Environement.TLSCertFile == "") != (db.Environement.TLSKeyFile == "") {
		return errors.New("TLS_CERT_FILE and TLS_KEY_FILE must be set together")
	}

	server := &http.Server{Addr: s.address, Handler: s.handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	errCh := make(chan error, 1)

	go func() {
		log.Println("Listening on", s.address)
		if db.Environement.TLSCertFile != "" {
			errCh <- server.ListenAndServeTLS(db.Environement.TLSCertFile, db.Environement.TLSKeyFile)
		} else {
			errCh <- server.ListenAndServe()
		}
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err

	case <-ctx.Done():
		// In case some requests are still being processed
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
