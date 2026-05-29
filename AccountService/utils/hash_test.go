package utils

import "testing"

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	const password = "correct horse battery staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hash == password {
		t.Fatal("HashPassword returned plaintext password")
	}
	if !CheckPasswordHash(password, hash) {
		t.Fatal("CheckPasswordHash rejected the original password")
	}
	if CheckPasswordHash("wrong password", hash) {
		t.Fatal("CheckPasswordHash accepted an invalid password")
	}
}
