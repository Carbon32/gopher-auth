package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	if len([]byte(password)) > 72 {
		return "", bcrypt.ErrPasswordTooLong
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hashed string, password []byte) bool {
	if len(password) > 72 {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashed), password)
	return err == nil
}
