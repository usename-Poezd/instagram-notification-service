package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/SmmTouch-com/instagram-notification-service/internal/config"
	"github.com/SmmTouch-com/instagram-notification-service/internal/pkg/logger"
	"github.com/SmmTouch-com/instagram-notification-service/pkg/mail"
	"go.uber.org/zap"
)

const ERROR_IDS_PATH = "logs/error_ids.json"

func main() {
	conf, err := config.Init(".env")
	if err != nil {
		log.Fatalf("Can't init config")
	}


	mailer := mail.NewMail(conf.MailConfig.Host, conf.MailConfig.Port, conf.MailConfig.From, conf.MailConfig.To, conf.MailConfig.Username, conf.MailConfig.Password)

	// Open our jsonFile
	jsonFile, err := os.Open(ERROR_IDS_PATH)
	// if we os.Open returns an error then handle it
	if err != nil {
		return
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

    var result []int
    json.Unmarshal([]byte(byteValue), &result)

	if len(result) == 0 {
		return
	}

	jsonByte, err := json.MarshalIndent(result, "", "    ")

	if err != nil {
		logger.GetLogger().Error("Can't marshal logs", zap.Error(err))
	}

	for _, v := range conf.MailConfig.To {
		if err := mailer.Send(v, "Ошибка рассылки Директ smmtouch.tech", "text/html", "<h3>Кол-во компаний:" + strconv.Itoa(len(result)) + "</h3><pre><h3>Id компаний:</h3><pre>"+string(jsonByte)+"</pre>"); err != nil {
			logger.GetLogger().Error("Error sending email", zap.String("email", v), zap.Error(err))
		}


		file, _ := json.MarshalIndent([]int{}, "", " ")
		_ = ioutil.WriteFile(ERROR_IDS_PATH, file, 0644)
	}

	defer jsonFile.Close()
}
