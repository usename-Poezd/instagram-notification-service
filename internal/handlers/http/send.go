package http

import (
	"fmt"
	"time"

	"github.com/SmmTouch-com/instagram-notification-service/internal/handlers/http/errors"
	"github.com/SmmTouch-com/instagram-notification-service/internal/pkg/logger"
	"github.com/SmmTouch-com/instagram-notification-service/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SendRequest struct {
	OrderId  int     `json:"order_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`

	Status   string `json:"status" validate:"required"`
	Currency string `json:"currency" validate:"required"`
	Tag      string `json:"tag" validate:"required"`
	Username string `json:"username" validate:"required"`
	Lang     string `json:"lang" validate:"required"`
}

var LastErrorTime time.Time

// Send
// @Summary Send message
// @Tags zeus
// @Description Send message
// @ModuleID ZeusAPI
// @Accept  json
// @Produce  json
// @Param input body SendRequest true "send data"
// @Success 200 {string} string "ok"
// @Failure 400,401,500,503 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ValidationErrorResponse
// @Router /send [post]
func (h *Handler) Send(c *fiber.Ctx) error {
	input := new(SendRequest)
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	var valErrors []errors.ValidationError
	if err := validator.New().Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errors.ValidationError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			valErrors = append(valErrors, element)
		}
	}

	if valErrors != nil {
		return c.Status(422).JSON(errors.ValidationErrorResponse{
			ErrorResponse: errors.ErrorResponse{
				Code:    422,
				Message: "Invalid request",
			},
			Errors: valErrors,
		})
	}

	id, err := h.services.Zeus.CreateCampaign(services.CampaignInput{
		OrderId:  input.OrderId,
		Status:   input.Status,
		Currency: input.Currency,
		Tag:      input.Tag,
		Quantity: input.Quantity,
		Username: input.Username,
		Amount:   input.Amount,
		Lang:     input.Lang,
	})

	if err != nil {
		logger.GetLogger().Error("Can't create Campaign", zap.Error(err))

		return c.Status(500).JSON(errors.ErrorResponse{
			Code:    500,
			Message: "Can't create Campaign",
		})
	}

	err = h.services.Zeus.StartCampaign(id)

	logger.GetLogger().Info("Campaign is started", zap.Int("id", id))

	if err != nil {
		logger.GetLogger().Error("Can't start Campaign", zap.Error(err))

		return c.Status(500).JSON(errors.ErrorResponse{
			Code:    500,
			Message: "Can't start Campaign",
		})
	}

	go func(id int) {
		time.Sleep(15 * time.Minute)

		isSent, err := h.services.Zeus.CheckMessageCampaignSent(id)
		if err != nil {
			logger.GetLogger().Error("Can't check message", zap.Error(err))
			fmt.Print(err)
		}

		logger.GetLogger().Info("Is message for Campaign sent?",
			zap.Int("id", id),
			zap.Bool("sent", isSent),
		)

		if isSent {
			err = h.services.Zeus.DeleteCampaign(id)

			if err != nil {
				logger.GetLogger().Error("Can't delete Campaign", zap.Error(err))
			}
		} else {
			logger.GetLogger().AddCompaignId(id)
		}
	}(id)
	return c.Status(200).SendString("ok")
}
