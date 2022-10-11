package app

import (
	"context"
	"log"

	"github.com/SmmTouch-com/instagram-notification-service/internal/config"
	"github.com/SmmTouch-com/instagram-notification-service/internal/handlers/http"
	"github.com/SmmTouch-com/instagram-notification-service/internal/services"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/googlesheets"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/mail"
	"github.com/gofiber/fiber/v2"
)

// @title Zeus Instagram DM API
// @version 1.0
// @description REST API for Zeus Instagram DM spam

// @host localhost:8000
// @BasePath /

// Run initializes whole application.
func Run() {
	app := fiber.New()

	conf, err := config.Init(".env")
	if err != nil {
		log.Fatalf("Can't init config")
	}

	ctx := context.Background()

	mailer := mail.NewMail(conf.MailConfig.Host, conf.MailConfig.Port, conf.MailConfig.From, conf.MailConfig.To, conf.MailConfig.Username, conf.MailConfig.Password)


	googleSheet := googlesheets.NewGoogleSheet(ctx, conf.LocaleConfig.Credentials, conf.LocaleConfig.SpreadSheetId)

	services := services.NewServices(services.Deps{
		Domain: conf.Domain,

		ZeusApiUrl: conf.ZeusConfig.ApiUrl,
		ZeusToken:  conf.ZeusConfig.Token,
		ZeusAdditionalRate: conf.ZeusConfig.AdditionalRate,

		MailTo: mailer.To,

		SheetKeyPrefix: conf.LocaleConfig.SheetKeyPrefix,
		SheetName:   conf.LocaleConfig.SheetName,
		GoogleSheet: googleSheet,
		Mailer: mailer,
	})

	handler := http.NewHandler(services)
	handler.Init(app)

	app.Listen(":8000")
}
