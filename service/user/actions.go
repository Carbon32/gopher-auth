package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"carb-backend/cmd/db"
	"carb-backend/service/auth"
	"carb-backend/service/mail"
	"carb-backend/service/types"
	"carb-backend/service/utils"
	"errors"
	"html"
)

func (h *Handler) issueAction(ctx context.Context, userID int, email, purpose, path, subject string) error {
	token, err := auth.RandomToken()
	if err != nil {
		return err
	}

	if err := h.store.CreateActionToken(ctx, userID, purpose, email, hashToken(token), time.Now().Add(30*time.Minute)); err != nil {
		return err
	}

	base := db.Environement.PublicAPIURL
	if purpose == "password_reset" {
		base = db.Environement.PublicAppURL
		path = "/reset_password"
	}

	link := strings.TrimRight(base, "/") + path + "?token=" + url.QueryEscape(token)
	return mail.SendEmail(email, subject, fmt.Sprintf("<p>Use this secure link to continue:</p><p><a href=\"%s\">Continue</a></p><p>This link expires in 30 minutes and can only be used once.</p>", link))
}

func (h *Handler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderActionForm(w, "Verify e-mail address", "Confirm e-mail", r.URL.Query().Get("token"), "/verify_email")
		return
	}

	var p struct {
		Token string `json:"token" validate:"required"`
	}

	if err := parseActionToken(w, r, &p.Token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	if err := h.store.CompleteEmailVerification(r.Context(), hashToken(p.Token)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid or expired verification link"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not verify account"))
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "E-mail verified"})
}

func (h *Handler) confirmEmailChange(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderActionForm(w, "Confirm e-mail change", "Confirm new e-mail", r.URL.Query().Get("token"), "/email/change/confirm")
		return
	}

	var p struct {
		Token string `json:"token" validate:"required"`
	}

	if err := parseActionToken(w, r, &p.Token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	if err := h.store.CompleteEmailChange(r.Context(), hashToken(p.Token)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid or expired e-mail confirmation link"))
		} else if isUniqueViolation(err) {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("E-mail is unavailable"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not update e-mail"))
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "E-mail address updated"})
}

func renderActionForm(w http.ResponseWriter, title, label, token, endpoint string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.WriteHeader(http.StatusOK)

	_, _ = fmt.Fprintf(w, `<!doctype html>
	<html lang="en">
	<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta name="color-scheme" content="dark">
	<title>%s</title>
	<style>
		:root { color-scheme: dark; font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #1e1e1e; color: #f1f1f1; }
		* { box-sizing: border-box; }
		body { min-height: 100vh; margin: 0; display: grid; place-items: center; padding: 24px; background: #1e1e1e; }
		main { width: min(100%%, 480px); padding: 38px; border: 1px solid #292929; border-radius: 14px; background: #131313; box-shadow: 0 18px 50px rgb(0 0 0 / 24%%); }
		.brand { margin: 0 0 34px; color: #fff; font-size: 13px; font-weight: 800; letter-spacing: .16em; }
		.eyebrow { margin: 0 0 10px; color: #666; font-size: 10px; font-weight: 700; letter-spacing: .16em; }
		h1 { margin: 0; color: #f1f1f1; font-size: clamp(25px, 7vw, 32px); font-weight: 600; letter-spacing: -.035em; line-height: 1.15; }
		.description { margin: 14px 0 26px; color: #888; font-size: 14px; line-height: 1.7; }
		button { width: 100%%; min-height: 44px; padding: 0 18px; border: 0; border-radius: 7px; background: #eee; color: #111; cursor: pointer; font: inherit; font-size: 13px; font-weight: 650; transition: background .2s ease, transform .2s ease; }
		button:hover { background: #fff; }
		button:active { transform: translateY(1px); }
		.footer { margin: 22px 0 0; color: #555; font-size: 11px; line-height: 1.6; }
		@media (max-width: 520px) { body { padding: 16px; } main { padding: 28px 24px; } .brand { margin-bottom: 28px; } }
		@media (prefers-reduced-motion: reduce) { button { transition: none; } }
	</style>
	</head>
	<body>
	<main>
		<p class="brand">SECURITY</p>
		<p class="eyebrow">ACCOUNT SECURITY</p>
		<h1>%s</h1>
		<p class="description">Confirm this one-time action to continue. For your security, this link expires after 30 minutes and can only be used once.</p>
		<form method="post" action="%s">
		<input type="hidden" name="token" value="%s">
		<button type="submit">%s</button>
		</form>
		<p class="footer">If you didn't request this, you can safely close this page.</p>
	</main>
	</body>
	</html>`,
		html.EscapeString(title),
		html.EscapeString(title),
		html.EscapeString("/api/v1"+endpoint),
		html.EscapeString(token),
		html.EscapeString(label),
	)
}

func parseActionToken(w http.ResponseWriter, r *http.Request, token *string) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var p struct {
			Token string `json:"token"`
		}

		if err := utils.ParseJSON(r, &p); err != nil {
			return err
		}

		*token = p.Token
		return nil
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		return err
	}

	*token = r.FormValue("token")
	return nil
}

func (h *Handler) resendVerification(w http.ResponseWriter, r *http.Request) {
	var p types.ForgotPassword
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	u, err := h.store.GetUserByEmail(p.Email)
	if err == nil && !u.Verified {
		if e := h.issueAction(r.Context(), u.Id, u.Email, "verify_email", "/verify_email", "Verify your email"); e != nil {
			log.Printf("verification e-mail error: %v", e)
		}
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]string{"message": "If verification is needed, an e-mail has been sent"})
}

func (h *Handler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var p types.ForgotPassword
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	p.Email = strings.ToLower(strings.TrimSpace(p.Email))
	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	u, err := h.store.GetUserByEmail(strings.TrimSpace(p.Email))
	if err == nil {
		if e := h.issueAction(r.Context(), u.Id, u.Email, "password_reset", "/password/reset", "Password reset"); e != nil {
			log.Printf("password reset e-mail error: %v", e)
		}
	} else if !errors.Is(err, ErrNotFound) {
		log.Printf("password reset lookup error: %v", err)
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]string{"message": "If that account exists, a password reset e-mail has been sent"})
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var p types.ResetPassword
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	hash, err := auth.HashPassword(p.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.CompletePasswordReset(r.Context(), hashToken(p.Token), hash); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid or expired reset link"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not update password"))
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password updated; please log in again"})
}

func (h *Handler) refreshSession(w http.ResponseWriter, r *http.Request) {
	var p struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid request"))
		return
	}

	newRefresh, err := auth.RandomToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	sid, uid, expires, err := h.store.RotateRefreshToken(r.Context(), hashToken(p.RefreshToken), hashToken(newRefresh))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Invalid or expired refresh token"))
		return
	}

	access, err := auth.CreateAccessJWT([]byte(db.Environement.JWTSecret), uid, sid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"access_token": access, "refresh_token": newRefresh, "expires_at": expires.UTC()})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Context().Value(auth.SessionIDKey).(string)
	if err := h.store.RevokeSession(r.Context(), sessionID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not log out"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}
