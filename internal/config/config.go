package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ZeusConfig
	LocaleConfig
	MailConfig

	Domain string `mapstructure:"DOMAIN"`
}

type ZeusConfig struct {
	Token  string `mapstructure:"ZEUS_TOKEN"`
	ApiUrl string `mapstructure:"ZEUS_API_URL"`
	AdditionalRate  float64 `mapstructure:"ZEUS_ADDITIONAL_RATE"`
}

type LocaleConfig struct {
	Credentials   string `mapstructure:"LOCALE_CREDENTIALS_FILE"`
	SpreadSheetId string `mapstructure:"LOCALE_SPREADSHEET_ID"`
	SheetName     string `mapstructure:"LOCALE_SHEET_NAME"`
	SheetKeyPrefix string `mapstructure:"LOCALE_KEY_PREFIX"`
}

type MailConfig struct {
	Port     int    `mapstructure:"MAIL_PORT"`
	Host     string `mapstructure:"MAIL_HOST"`
	From     string `mapstructure:"MAIL_FROM"`
	Username string `mapstructure:"MAIL_USERNAME"`
	Password string `mapstructure:"MAIL_PASSWORD"`

	To []string
}

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configFile string) (*Config, error) {

	if err := parseConfigFile(configFile); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigFile(configFile string) error {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}
		

	err = viper.Unmarshal(&cfg.LocaleConfig)
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&cfg.MailConfig)
	cfg.MailConfig.To = strings.Split(viper.GetString("MAIL_TO"), ",")
	if err != nil {
		return err
	}

	return viper.Unmarshal(&cfg.ZeusConfig)
}
