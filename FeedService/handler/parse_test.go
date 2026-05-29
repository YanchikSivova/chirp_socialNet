package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestParseFeedIdParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing id", func(t *testing.T) {
		ctx := newTestContext("/")

		if id, err := parseIdParam(ctx); err == nil || id != nil {
			t.Fatalf("parseIdParam() = %v, %v; want nil id and error", id, err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		ctx := newTestContext("/")
		ctx.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

		if id, err := parseIdParam(ctx); err == nil || id != nil {
			t.Fatalf("parseIdParam() = %v, %v; want nil id and error", id, err)
		}
	})

	t.Run("valid id", func(t *testing.T) {
		want := uuid.New()
		ctx := newTestContext("/")
		ctx.Params = gin.Params{{Key: "id", Value: want.String()}}

		got, err := parseIdParam(ctx)
		if err != nil {
			t.Fatalf("parseIdParam returned error: %v", err)
		}
		if got == nil || *got != want {
			t.Fatalf("id = %v, want %s", got, want)
		}
	})
}

func newTestContext(target string) *gin.Context {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, target, nil)
	return ctx
}
