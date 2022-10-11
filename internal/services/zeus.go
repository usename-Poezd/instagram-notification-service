package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SmmTouch-com/instagram-notification-service/internal/domain"
	"github.com/SmmTouch-com/instagram-notification-service/internal/pkg/logger"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/mail"
	"go.uber.org/zap"
)

type ZeusService struct {
	Domain string
	Token  string
	ApiUrl string
	mailTo []string
	additionalRate float64
	Locale Locale
	mailer mail.Mailer
}

func NewZeusService(Domain string, Token string, ApiUrl string, MailTo []string, AdditionalRate float64, Locale Locale, Mailer mail.Mailer) *ZeusService {
	return &ZeusService{
		Domain,
		Token,
		ApiUrl,
		MailTo,
		AdditionalRate,
		Locale,
		Mailer,
	}
}

func (s *ZeusService) makeRequest(method string, url string, input any, unmarshalValue any) error {

	reqBody := []byte("")

	if input != nil {
		reqBody, _ = json.Marshal(input)
	}

	client := http.Client{}
	req, _ := http.NewRequest(method, s.ApiUrl+url, bytes.NewReader(reqBody))

	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		for _, v := range s.mailTo {
			if err := s.mailer.Send(v,
				"Ошибка рассылки Директ smmtouch.tech ZEUS API ERROR", "text/html",
				"<p>Url: <b>"+s.ApiUrl+url+"</b>"+
					"<p>Метод: <b>"+method+"</b></p><pre>"+err.Error()+"</pre>"); err != nil {
				logger.GetLogger().Error("Error sending email", zap.String("email", v), zap.Error(err))
			}
		}

		logger.GetLogger().Error("Error of sending request", zap.Error(err))
		return err
	}

	if res.StatusCode != http.StatusOK {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.GetLogger().Error("client: could not read response body", zap.Error(err))
			return err
		}

		var result map[string]interface{}
		err = json.Unmarshal(resBody, &result)
		if err != nil {
			logger.GetLogger().Error("client: could not unmarshal response body", zap.Error(err))
			return err
		}

		jsonByte, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			logger.GetLogger().Error("client: could not marshal response body", zap.Error(err))
			return err
		}

		for _, v := range s.mailTo {
			if err := s.mailer.Send(v,
				"Ошибка рассылки Директ smmtouch.tech ZEUS API ERROR", "text/html",
				"<p>Ошибка со статусом: <b>"+strconv.Itoa(res.StatusCode)+
					"</b></p><p>Url: <b>"+s.ApiUrl+url+"</b>"+
					"<p>Метод: <b>"+req.Method+"</b></p>"+"<pre>"+string(jsonByte)+"</pre>"); err != nil {
				logger.GetLogger().Error("Error sending email", zap.String("email", v), zap.Error(err))
			}
		}

		logger.GetLogger().Error("Zeus API request bad satatus", zap.Any("response", result))

		return errors.New("Zeus API request bad satatus")
	}

	if unmarshalValue != nil {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.GetLogger().Error("client: could not read response body", zap.Error(err))
			return err
		}

		err = json.Unmarshal(resBody, unmarshalValue)
		if err != nil {
			logger.GetLogger().Error("client: could not unmarshal response body", zap.Error(err))
			return err
		}
	}

	return nil
}

type CampaignMessagesResponse struct {
	Messages []struct {
		TagetUserName string `json:"target_user_name"`
	} `json:"messages"`
}

func (s *ZeusService) GetCampaign(id int) (*domain.Campaign, error) {
	logger.GetLogger().Info("Getting Campaign", zap.Int("id", id))

	resCampaign := new(domain.Campaign)

	err := s.makeRequest(http.MethodGet, "direct/cpa/"+strconv.Itoa(id)+"/", nil, resCampaign)
	if err != nil {
		return nil, err
	}

	return resCampaign, nil
}

func (s *ZeusService) CheckMessageCampaignSent(id int) (bool, error) {
	logger.GetLogger().Info("Checking is message sent for Campaign", zap.Int("id", id))

	campaignMessagesResponse := new(CampaignMessagesResponse)

	err := s.makeRequest(http.MethodGet, "direct/cpa/"+strconv.Itoa(id)+"/messages?page=0", nil, campaignMessagesResponse)
	if err != nil {
		return false, err
	}

	return len(campaignMessagesResponse.Messages) != 0, nil
}

func (s *ZeusService) CreateCampaign(input CampaignInput) (int, error) {
	logger.GetLogger().Info("Creating Campaign", zap.Any("input", input))

	messageCh := make(chan string)
	tagCh := make(chan string)

	errCh := make(chan error, 2)

	go func() {
		messageVal, err := s.Locale.GetMessageByKeyWithPrefix(input.Status, input.Lang)
		if err != nil {
			errCh <- err
		}

		errCh <- nil
		messageCh <- messageVal
	}()

	go func() {
		tagVal, err := s.Locale.GetMessageByKey(input.Tag, input.Lang)
		if err != nil {
			errCh <- err
		}

		errCh <- nil
		tagCh <- tagVal
	}()

	err := <-errCh

	if err != nil {
		return 0, err
	}

	message := <-messageCh
	tag := <-tagCh
	message = strings.ReplaceAll(message, "{tag}", tag)
	message = strings.ReplaceAll(message, "{username}", input.Username)
	message = strings.ReplaceAll(message, "{quantity}", strconv.Itoa(input.Quantity))
	message = strings.ReplaceAll(message, "{amount}", strconv.FormatFloat(input.Amount, 'f', -1, 64)+" "+input.Currency)

	Campaign := domain.Campaign{
		DayStartTime:     0,
		DayEndTime:       0,
		WorkRounTheClock: true,
		Name:             input.Username + "_" + s.Domain + "_" + input.Status,
		Id:               0,
		Started:          true,
		UsaIgnoreList:    false,
		Message:          message,
		MessageType:      "TEXT",
		Type:             "ON_LIST",
		UserList:         input.Username,
		AdditionalRate:   s.additionalRate,
	}

	resCampaign := new(domain.Campaign)

	logger.GetLogger().Info("Request for creating Campaign", zap.Any("Campaign", Campaign))
	err = s.makeRequest(http.MethodPost, "direct/cpa/", Campaign, resCampaign)
	if err != nil {
		return 0, err
	}

	return resCampaign.Id, nil
}

func (s *ZeusService) DeleteCampaign(id int) error {
	logger.GetLogger().Info("Deleting Campaign", zap.Int("id", id))
	return s.makeRequest(http.MethodDelete, "direct/cpa/"+strconv.Itoa(id)+"/", nil, nil)
}

func (s *ZeusService) StartCampaign(id int) error {

	retry := 3

	for i := 1; i <= retry; i++ {
		logger.GetLogger().Info("Starting Campaign id try:"+strconv.Itoa(i), zap.Int("id", id))
		err := s.makeRequest(http.MethodPost, "direct/cpa/"+strconv.Itoa(id)+"/start", nil, nil)
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Millisecond)

		campign, err := s.GetCampaign(id)
		if err != nil {
			return err
		}

		if campign.Started {
			return nil
		}
	}

	logger.GetLogger().Info("Can't start champagin", zap.Int("id", id))

	return errors.New("Can't start champagin")
}

func (s *ZeusService) StopCampaign(id int) error {
	logger.GetLogger().Info("Stopping Campaign id", zap.Int("id", id))
	return s.makeRequest(http.MethodPost, "direct/cpa/"+strconv.Itoa(id)+"/stop", nil, nil)
}

type CampaignLogs struct {
	Total int `json:"total"`

	Logs []struct {
		ActionTime  int    `json:"action_time"`
		Description string `json:"description"`
		LogLevel    string `json:"log_level"`
		ErrorType   string `json:"error_type"`
	}
}

func (s *ZeusService) GetLogs(id int) (interface{}, error) {
	logger.GetLogger().Info("Getting Campaign logs", zap.Int("id", id))

	res := new(CampaignLogs)

	err := s.makeRequest(http.MethodGet, "direct/cpa/"+strconv.Itoa(id)+"/logs?page=0&pageSize=50", nil, res)
	if err != nil {
		return res, err
	}

	return res, nil
}
