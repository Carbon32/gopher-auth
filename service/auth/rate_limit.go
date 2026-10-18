package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Both local to this file

type bucket struct {
	start time.Time
	count int
}
type LoginLimiter struct {
	mu      sync.Mutex
	buckets map[string]bucket
	limit   int
	window  time.Duration
}

func NewLoginLimiter(limit int, window time.Duration) *LoginLimiter {
	return &LoginLimiter{buckets: make(map[string]bucket), limit: limit, window: window}
}

func (l *LoginLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		host = strings.TrimSpace(host)
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20)) // If you wanna be smart
		if err != nil {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		var fields struct {
			Email    string `json:"email"`
			Username string `json:"username"`
		}
		_ = json.Unmarshal(body, &fields)
		identity := strings.ToLower(strings.TrimSpace(fields.Email))
		if identity == "" {
			identity = strings.ToLower(strings.TrimSpace(fields.Username))
		}
		keys := []string{"ip:" + host}
		if identity != "" {
			sum := sha256.Sum256([]byte(identity))
			keys = append(keys, "account:"+hex.EncodeToString(sum[:]))
		}
		now := time.Now()
		l.mu.Lock()
		for key, b := range l.buckets {
			if now.Sub(b.start) >= l.window {
				delete(l.buckets, key)
			}
		}
		if len(l.buckets) >= 50000 && len(keys) > 1 {
			keys = keys[:1]
		}
		for _, key := range keys {
			if b := l.buckets[key]; b.count >= l.limit && now.Sub(b.start) < l.window {
				l.mu.Unlock()
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too many login attempts", http.StatusTooManyRequests)
				return
			}
		}
		for _, key := range keys {
			b := l.buckets[key]
			if b.start.IsZero() || now.Sub(b.start) >= l.window {
				b = bucket{start: now}
			}
			b.count++
			l.buckets[key] = b
		}
		l.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
