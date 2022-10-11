package http

import (
	_ "github.com/SmmTouch-com/instagram-notification-service/docs"
	"github.com/SmmTouch-com/instagram-notification-service/internal/services"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{
		services,
	}
}

func (h Handler) Init(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/ping", h.Ping)

	app.Post("/send", h.Send)
}

// Ping
// @Summary Ping
// @Tags service
// @Description Ping
// @ModuleID Зштп
// @Accept  json
// @Produce  json
// @Success 200 {string} string "pong"
// @Failure 400,401,500,503 {string} string "error"
// @Router /ping [get]
func (h Handler) Ping(c *fiber.Ctx) error {
	return c.Status(200).SendString("pong")
}
