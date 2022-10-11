package services

import (
	"github.com/SmmTouch-com/instagram-notification-service/internal/pkg/logger"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/googlesheets"
	"go.uber.org/zap"
)

type LocaleService struct {
	googleSheet googlesheets.Sheet
	sheetName   string
	prefix string
}

func NewLocaleService(sheetName string, prefix string, googleSheet googlesheets.Sheet) *LocaleService {
	return &LocaleService{
		googleSheet,
		sheetName,
		prefix,
	}
}

func (s *LocaleService) GetMessageByKeyWithPrefix(key string, lang string) (string, error) {
	return s.GetMessageByKey(s.prefix + key, lang)
}

func (s *LocaleService) GetMessageByKey(key string, lang string) (string, error) {

	logger.GetLogger().Info("Getting message", zap.String("key", key), zap.String("lang", lang))

	rang, err := s.googleSheet.Get(s.sheetName)
	if err != nil {
		return "", err
	}

	langIdx := -1

	for i, value := range rang.Values[0] {
		if lang == value {
			langIdx = i
		}
	}

	if langIdx < 0 {
		return "", nil
	}

	for _, value := range rang.Values {
		if len(value) > 1 && value[0] == key {
			return value[langIdx].(string), nil
		}
	}

	return "", nil
}
