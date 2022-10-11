package services

import (
	"github.com/SmmTouch-com/instagram-notification-service/internal/domain"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/googlesheets"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/mail"
)

type Deps struct {
	Domain string

	ZeusApiUrl         string
	ZeusToken          string
	ZeusAdditionalRate float64

	MailTo []string

	SheetName      string
	SheetKeyPrefix string
	GoogleSheet    googlesheets.Sheet
	Mailer         mail.Mailer
}

type Services struct {
	Zeus   Zeus
	Locale Locale
}

type CampaignInput struct {
	OrderId  int
	Quantity int
	Amount   float64
	Status   string
	Currency string
	Tag      string
	Username string
	Lang     string
}

type Zeus interface {
	GetCampaign(id int) (*domain.Campaign, error)
	CreateCampaign(input CampaignInput) (int, error)
	DeleteCampaign(id int) error

	StartCampaign(id int) error
	StopCampaign(id int) error

	CheckMessageCampaignSent(id int) (bool, error)
	GetLogs(id int) (interface{}, error)
}

type Locale interface {
	GetMessageByKey(key string, lang string) (string, error)
	GetMessageByKeyWithPrefix(key string, lang string) (string, error)
}

func NewServices(deps Deps) *Services {
	Locale := NewLocaleService(deps.SheetName, deps.SheetKeyPrefix, deps.GoogleSheet)

	return &Services{
		Zeus:   NewZeusService(deps.Domain, deps.ZeusToken, deps.ZeusApiUrl, deps.MailTo, deps.ZeusAdditionalRate, Locale, deps.Mailer),
		Locale: Locale,
	}
}
