package services_test

import (
	"errors"

	"github.com/SmmTouch-com/instagram-notification-service/internal/services"
	"github.com/golang/mock/gomock"
	"google.golang.org/api/sheets/v4"
)

func (s *ServicesTestSuite) TestNewZeusService_CreateCampaignLocaleError() {
	localeService := services.NewLocaleService("test", "prefix", s.mocks.sheet)
	zeusService := services.NewZeusService("example.com", "token", "apiUrl", []string{}, 0, localeService, s.mocks.mailer)

	s.mocks.sheet.EXPECT().Get(gomock.Any()).Return(&sheets.ValueRange{}, errInternalServErr)

	res, err := zeusService.CreateCampaign(services.CampaignInput{})

	s.True(errors.Is(err, errInternalServErr))
	s.Equal(0, res)
}
