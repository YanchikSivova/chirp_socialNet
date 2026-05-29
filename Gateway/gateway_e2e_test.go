package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"gateway/middleware"
	"gateway/proxy"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestGatewayE2EForwardsAuthenticatedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "gateway-secret")

	const profileID = "3f649716-ccf1-41fc-8760-43dc0d855b11"
	var gotProfileID string
	var gotPath string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotProfileID = r.Header.Get("X-User-Id")
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(gin.H{"success": true})
	}))
	defer upstream.Close()

	reverseProxy, err := proxy.NewReverseProxy(upstream.URL)
	if err != nil {
		t.Fatalf("NewReverseProxy returned error: %v", err)
	}

	router := gin.New()
	protectedUsers := router.Group("/users")
	protectedUsers.Use(middleware.AuthMiddleware())
	protectedUsers.GET("/me", proxy.ProxyHandler(reverseProxy))

	gateway := httptest.NewServer(router)
	defer gateway.Close()

	req, err := http.NewRequest(http.MethodGet, gateway.URL+"/users/me", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+signGatewayToken(t, profileID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("gateway request returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if gotProfileID != profileID {
		t.Fatalf("X-User-Id = %q, want %q", gotProfileID, profileID)
	}
	if gotPath != "/users/me" {
		t.Fatalf("path = %q, want /users/me", gotPath)
	}
}

func TestGatewayE2ERejectsMissingTokenBeforeProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "gateway-secret")

	var upstreamWasCalled atomic.Bool
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamWasCalled.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	reverseProxy, err := proxy.NewReverseProxy(upstream.URL)
	if err != nil {
		t.Fatalf("NewReverseProxy returned error: %v", err)
	}

	router := gin.New()
	protectedUsers := router.Group("/users")
	protectedUsers.Use(middleware.AuthMiddleware())
	protectedUsers.GET("/me", proxy.ProxyHandler(reverseProxy))

	gateway := httptest.NewServer(router)
	defer gateway.Close()

	resp, err := http.Get(gateway.URL + "/users/me")
	if err != nil {
		t.Fatalf("gateway request returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
	if upstreamWasCalled.Load() {
		t.Fatal("upstream was called for unauthenticated request")
	}
}

func signGatewayToken(t *testing.T, profileID string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"profile_id": profileID,
		"exp":        time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("gateway-secret"))
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}
	return tokenStr
}
