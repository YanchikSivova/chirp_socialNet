package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestDockerE2ELoginAndGetCurrentUser(t *testing.T) {
	if os.Getenv("RUN_DOCKER_E2E") != "1" {
		t.Skip("set RUN_DOCKER_E2E=1 to run this test against docker compose")
	}

	baseURL := getenv("E2E_BASE_URL", "http://localhost:8000")
	email := getenv("E2E_EMAIL", "test@email.com")
	password := getenv("E2E_PASSWORD", "password")

	client := &http.Client{Timeout: 5 * time.Second}
	loginBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		t.Fatalf("Marshal login body returned error: %v", err)
	}

	loginResp, err := client.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("POST /auth/login returned error: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("POST /auth/login status = %d, want %d", loginResp.StatusCode, http.StatusOK)
	}

	var loginPayload struct {
		AccessToken string `json:"access_token"`
	}
	if err = json.NewDecoder(loginResp.Body).Decode(&loginPayload); err != nil {
		t.Fatalf("Decode login response returned error: %v", err)
	}
	if loginPayload.AccessToken == "" {
		t.Fatal("login response does not contain access_token")
	}

	meReq, err := http.NewRequest(http.MethodGet, baseURL+"/users/me", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	meReq.Header.Set("Authorization", "Bearer "+loginPayload.AccessToken)

	meResp, err := client.Do(meReq)
	if err != nil {
		t.Fatalf("GET /users/me returned error: %v", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/me status = %d, want %d", meResp.StatusCode, http.StatusOK)
	}

	var mePayload struct {
		ProfileID string `json:"profile_id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
	}
	if err = json.NewDecoder(meResp.Body).Decode(&mePayload); err != nil {
		t.Fatalf("Decode /users/me response returned error: %v", err)
	}
	if mePayload.ProfileID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("profile_id = %q, want seeded profile id", mePayload.ProfileID)
	}
	if mePayload.Name != "Test User" {
		t.Fatalf("name = %q, want Test User", mePayload.Name)
	}
	if mePayload.Username != "test_user" {
		t.Fatalf("username = %q, want test_user", mePayload.Username)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
