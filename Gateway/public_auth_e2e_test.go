package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gateway/proxy"

	"github.com/gin-gonic/gin"
)

func TestGatewayE2EPublicAuthLoginProxiesRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotMethod string
	var gotPath string
	var gotBody map[string]string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode request body returned error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(gin.H{"access_token": "token-from-upstream"})
	}))
	defer upstream.Close()

	reverseProxy, err := proxy.NewReverseProxy(upstream.URL)
	if err != nil {
		t.Fatalf("NewReverseProxy returned error: %v", err)
	}

	router := gin.New()
	publicAuth := router.Group("/auth")
	publicAuth.POST("/login", proxy.ProxyHandler(reverseProxy))

	gateway := httptest.NewServer(router)
	defer gateway.Close()

	reqBody := `{"email":"user@example.com","password":"secret"}`
	resp, err := http.Post(gateway.URL+"/auth/login", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("gateway request returned error: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Decode response body returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q, want %q", gotMethod, http.MethodPost)
	}
	if gotPath != "/auth/login" {
		t.Fatalf("path = %q, want /auth/login", gotPath)
	}
	if gotBody["email"] != "user@example.com" || gotBody["password"] != "secret" {
		t.Fatalf("body = %#v, want login credentials", gotBody)
	}
	if response["access_token"] != "token-from-upstream" {
		t.Fatalf("access_token = %q, want token-from-upstream", response["access_token"])
	}
}
