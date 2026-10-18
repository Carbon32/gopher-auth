package user

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"carb-backend/service/types"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = sql.ErrNoRows

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(`INSERT INTO users (first_name,last_name,username,email,password,verified) VALUES ($1,$2,$3,$4,$5,$6)`, user.FirstName, user.LastName, user.Username, user.Email, user.Password, user.Verified)
	return err
}

func (s *Store) UpdateUser(user types.User) error {
	_, err := s.db.Exec(`UPDATE users SET first_name=$1,last_name=$2,username=$3,email=$4 WHERE id=$5`, user.FirstName, user.LastName, user.Username, user.Email, user.Id)
	return err
}

func (s *Store) ChangePassword(ctx context.Context, id int, passwordHash string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `UPDATE users SET password=$1 WHERE id=$2`, passwordHash, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE auth_sessions SET revoked_at=now() WHERE user_id=$1 AND revoked_at IS NULL`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) CompleteEmailVerification(ctx context.Context, tokenHash string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int
	if err = tx.QueryRowContext(ctx, `UPDATE auth_action_tokens SET consumed_at=now() WHERE token_hash=$1 AND purpose='verify_email' AND consumed_at IS NULL AND expires_at>now() RETURNING user_id`, tokenHash).Scan(&id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE users SET verified=TRUE WHERE id=$1`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) CompleteEmailChange(ctx context.Context, tokenHash string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int
	var email sql.NullString
	if err = tx.QueryRowContext(ctx, `UPDATE auth_action_tokens SET consumed_at=now() WHERE token_hash=$1 AND purpose='email_change' AND consumed_at IS NULL AND expires_at>now() RETURNING user_id,target_email`, tokenHash).Scan(&id, &email); err != nil {
		return err
	}
	if !email.Valid || email.String == "" {
		return sql.ErrNoRows
	}
	if _, err = tx.ExecContext(ctx, `UPDATE users SET email=$1,verified=TRUE WHERE id=$2`, email.String, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) CompletePasswordReset(ctx context.Context, tokenHash, passwordHash string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int
	if err = tx.QueryRowContext(ctx, `UPDATE auth_action_tokens SET consumed_at=now() WHERE token_hash=$1 AND purpose='password_reset' AND consumed_at IS NULL AND expires_at>now() RETURNING user_id`, tokenHash).Scan(&id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE users SET password=$1 WHERE id=$2`, passwordHash, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE auth_sessions SET revoked_at=now() WHERE user_id=$1 AND revoked_at IS NULL`, id); err != nil {
		return err
	}
	return tx.Commit()
}

const userColumns = `id, first_name, last_name, username, email, password, created_at, verified`

func (s *Store) GetUserById(id int) (*types.User, error) {
	return scanUser(s.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id=$1`, id))
}
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	return scanUser(s.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE lower(email)=lower($1)`, email))
}
func (s *Store) GetUserByUsername(username string) (*types.User, error) {
	return scanUser(s.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE lower(username)=lower($1)`, username))
}

func scanUser(row *sql.Row) (*types.User, error) {
	u := new(types.User)
	err := row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.Verified)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) CreateActionToken(ctx context.Context, userID int, purpose, targetEmail, tokenHash string, expiresAt time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `DELETE FROM auth_action_tokens WHERE consumed_at IS NOT NULL OR expires_at<=now()`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM auth_action_tokens WHERE user_id=$1 AND purpose=$2 AND consumed_at IS NULL`, userID, purpose); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO auth_action_tokens(token_hash,user_id,purpose,target_email,expires_at) VALUES($1,$2,$3,$4,$5)`, tokenHash, userID, purpose, targetEmail, expiresAt); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) CreateSession(ctx context.Context, id string, userID int, refreshHash string, expiresAt time.Time) error {
	_, _ = s.db.ExecContext(ctx, `DELETE FROM auth_sessions WHERE revoked_at IS NOT NULL OR expires_at<=now()`)
	_, err := s.db.ExecContext(ctx, `INSERT INTO auth_sessions(id,user_id,refresh_hash,expires_at) VALUES($1,$2,$3,$4)`, id, userID, refreshHash, expiresAt)
	return err
}
func (s *Store) SessionActive(ctx context.Context, id string) (bool, error) {
	var active bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth_sessions WHERE id=$1 AND revoked_at IS NULL AND expires_at>now())`, id).Scan(&active)
	return active, err
}
func (s *Store) RotateRefreshToken(ctx context.Context, oldRefreshHash, newRefreshHash string) (string, int, time.Time, error) {
	var id string
	var userID int
	var expires time.Time
	err := s.db.QueryRowContext(ctx, `UPDATE auth_sessions SET refresh_hash=$1 WHERE refresh_hash=$2 AND revoked_at IS NULL AND expires_at>now() RETURNING id,user_id,expires_at`, newRefreshHash, oldRefreshHash).Scan(&id, &userID, &expires)
	return id, userID, expires, err
}
func (s *Store) RevokeSession(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE auth_sessions SET revoked_at=now() WHERE id=$1 AND revoked_at IS NULL`, id)
	return err
}
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
