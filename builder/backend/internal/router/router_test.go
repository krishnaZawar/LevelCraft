package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func Test_CreateRoutes(t *testing.T) {
	app := fiber.New()

	CreateRoutes(app)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /ping status = %d, want = %d", resp.StatusCode, http.StatusOK)
	}
}
