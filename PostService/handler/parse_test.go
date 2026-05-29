package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestParseProfileHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing header", func(t *testing.T) {
		ctx := newTestContext("/")

		if id, err := parseProfileHeader(ctx); err == nil || id != nil {
			t.Fatalf("parseProfileHeader() = %v, %v; want nil id and error", id, err)
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		ctx := newTestContext("/")
		ctx.Request.Header.Set("X-User-Id", "not-a-uuid")

		if id, err := parseProfileHeader(ctx); err == nil || id != nil {
			t.Fatalf("parseProfileHeader() = %v, %v; want nil id and error", id, err)
		}
	})

	t.Run("valid uuid", func(t *testing.T) {
		want := uuid.New()
		ctx := newTestContext("/")
		ctx.Request.Header.Set("X-User-Id", want.String())

		got, err := parseProfileHeader(ctx)
		if err != nil {
			t.Fatalf("parseProfileHeader returned error: %v", err)
		}
		if got == nil || *got != want {
			t.Fatalf("profile id = %v, want %s", got, want)
		}
	})
}

func TestParseIdParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing id", func(t *testing.T) {
		ctx := newTestContext("/")

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

func TestParsePaginationQueries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("defaults", func(t *testing.T) {
		ctx := newTestContext("/")

		offset, err := parseOffsetQuery(ctx)
		if err != nil {
			t.Fatalf("parseOffsetQuery returned error: %v", err)
		}
		if offset == nil || *offset != 0 {
			t.Fatalf("offset = %v, want 0", offset)
		}

		limit, err := parseLimitQuery(ctx)
		if err != nil {
			t.Fatalf("parseLimitQuery returned error: %v", err)
		}
		if limit == nil || *limit != 20 {
			t.Fatalf("limit = %v, want 20", limit)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		ctx := newTestContext("/?offset=7&limit=12")

		offset, err := parseOffsetQuery(ctx)
		if err != nil {
			t.Fatalf("parseOffsetQuery returned error: %v", err)
		}
		if offset == nil || *offset != 7 {
			t.Fatalf("offset = %v, want 7", offset)
		}

		limit, err := parseLimitQuery(ctx)
		if err != nil {
			t.Fatalf("parseLimitQuery returned error: %v", err)
		}
		if limit == nil || *limit != 12 {
			t.Fatalf("limit = %v, want 12", limit)
		}
	})

	t.Run("invalid values", func(t *testing.T) {
		ctx := newTestContext("/?offset=bad&limit=bad")

		if offset, err := parseOffsetQuery(ctx); err == nil || offset != nil {
			t.Fatalf("parseOffsetQuery() = %v, %v; want nil offset and error", offset, err)
		}
		if limit, err := parseLimitQuery(ctx); err == nil || limit != nil {
			t.Fatalf("parseLimitQuery() = %v, %v; want nil limit and error", limit, err)
		}
	})
}

func newTestContext(target string) *gin.Context {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, target, nil)
	return ctx
}
