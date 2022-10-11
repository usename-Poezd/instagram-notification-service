package http_test

import (
	"net/http/httptest"
	"testing"

	handler "github.com/SmmTouch-com/instagram-notification-service/internal/handlers/http"
	"github.com/SmmTouch-com/instagram-notification-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	h := handler.NewHandler(&services.Services{})

	require.IsType(t, &handler.Handler{}, h)
}

func TestNewHandler_Routing(t *testing.T) {
	app := fiber.New()

	h := handler.NewHandler(&services.Services{})
	h.Init(app)

	req := httptest.NewRequest("GET", "/ping", nil)
	

	res, _ := app.Test(req, -1)

	assert.Equalf(t, res.StatusCode, 200, "Test ping")
}