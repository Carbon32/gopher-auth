package auth

import "testing"

func TestPasswordHashAndCompare(t *testing.T) {
	password := "an example password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if !ComparePassword(hash, []byte(password)) {
		t.Fatal("correct password did not match")
	}
	if ComparePassword(hash, []byte("a different password")) {
		t.Fatal("incorrect password matched")
	}
}

func TestPasswordOverBcryptLimitRejected(t *testing.T) {
	if _, err := HashPassword(string(make([]byte, 73))); err == nil {
		t.Fatal("expected bcrypt byte limit error")
	}
}
