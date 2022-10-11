package services_test

import (
	"errors"

	"github.com/SmmTouch-com/instagram-notification-service/internal/services"
	"github.com/golang/mock/gomock"
	"google.golang.org/api/sheets/v4"
)

var errInternalServErr = errors.New("test: internal server error")



func (s *ServicesTestSuite) TestNewLocaleService_GetMessageByKey() {
	localeService := services.NewLocaleService("test", "prefix", s.mocks.sheet)

	s.mocks.sheet.EXPECT().Get(gomock.Any()).Return(
		&sheets.ValueRange{
			Values: [][]interface{}{
				{"key", "ru", "en"},
				{"testKeyAnother", "ruName1", "enName1"},
				{"testKey", "ruName", "enName"},
			},
		},
		nil)

	res, err := localeService.GetMessageByKey("testKey", "ru")

	s.False(errors.Is(err, errInternalServErr))
	s.Equal("ruName", res)
}

func (s *ServicesTestSuite) TestNewLocaleService_GetMessageByKeyWithPrefix() {
	localeService := services.NewLocaleService("test", "prefix.", s.mocks.sheet)

	s.mocks.sheet.EXPECT().Get(gomock.Any()).Return(
		&sheets.ValueRange{
			Values: [][]interface{}{
				{"key", "ru", "en"},
				{"testKeyAnother", "ruName1", "enName1"},
				{"prefix.testKey", "ruName", "enName"},
			},
		},
		nil)

	res, err := localeService.GetMessageByKeyWithPrefix("testKey", "ru")

	s.False(errors.Is(err, errInternalServErr))
	s.Equal("ruName", res)
}

func (s *ServicesTestSuite) TestNewLocaleService_GetMessageByKeyNotExistsLang() {
	localeService := services.NewLocaleService("test", "prefix", s.mocks.sheet)

	s.mocks.sheet.EXPECT().Get(gomock.Any()).Return(
		&sheets.ValueRange{
			Values: [][]interface{}{
				{"key", "ru", "en"},
				{"testKeyAnother", "ruName1", "enName1"},
				{"testKey", "ruName", "enName"},
			},
		},
		nil)

	res, err := localeService.GetMessageByKey("testKey", "pt")

	s.Equal(nil, err)
	s.Equal("", res)
}

func (s *ServicesTestSuite) TestNewLocaleService_GetMessageByKeyNotExistsKey() {
	localeService := services.NewLocaleService("test", "prefix", s.mocks.sheet)

	s.mocks.sheet.EXPECT().Get(gomock.Any()).Return(
		&sheets.ValueRange{
			Values: [][]interface{}{
				{"key", "ru", "en"},
				{"testKeyAnother", "ruName1", "enName1"},
				{"testKey", "ruName", "enName"},
			},
		},
		nil)

	res, err := localeService.GetMessageByKey("testKeyNotExists", "ru")

	s.Equal(nil, err)
	s.Equal("", res)
}


func (s *ServicesTestSuite) TestNewLocaleService_GetMessageByKeyErr() {
	localeService := services.NewLocaleService("test", "prefix", s.mocks.sheet)

	s.mocks.sheet.EXPECT().Get(gomock.Any()).Return(&sheets.ValueRange{}, errInternalServErr)

	res, err := localeService.GetMessageByKey("key", "string")

	s.True(errors.Is(err, errInternalServErr))
	s.Equal("", res)
}