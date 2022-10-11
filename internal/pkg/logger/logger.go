package logger

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lock = &sync.Mutex{}

type Logger struct {
	log *zap.Logger
}

var logetInstance *Logger

func NewLogger() *Logger {
	conf := zap.NewProductionConfig()

	conf.OutputPaths = []string{
		"stdout",
		"./logs/service.log",
	}

	if flag.Lookup("test.v") != nil {
		conf.OutputPaths = []string{
			"stdout",
		}
	}

	logger, err := conf.Build()
	if err != nil {
		log.Fatal(err)
	}

	return &Logger{
		log: logger,
	}
}

func GetLogger() *Logger {
	if logetInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if logetInstance == nil {

			logetInstance = NewLogger()
		}
	}

	return logetInstance
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.log.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.log.Error(msg, fields...)
}

const ERROR_IDS_PATH = "logs/error_ids.json"
func (l *Logger) AddCompaignId(id int) {
	// Open our jsonFile
	jsonFile, err := os.Open(ERROR_IDS_PATH)
	// if we os.Open returns an error then handle it
	if err != nil {
		file, _ := json.MarshalIndent([]int{id}, "", " ")
		_ = ioutil.WriteFile(ERROR_IDS_PATH, file, 0644)

		return
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

    var result []int
    json.Unmarshal([]byte(byteValue), &result)

	result = append(result, id)

	file, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile(ERROR_IDS_PATH, file, 0644)
}
