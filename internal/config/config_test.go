package config

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestInit(t *testing.T) {
	type env struct {
		Domain string

		Token  string
		ApiUrl string
		AdditionalRate string

		Credentials   string
		SpreadSheetId string
		SheetName     string
		SheetKeyPrefix string

		Port int
		Host string
		From string
		Username string
		Password string
		To string
	}

	type args struct {
		path string
		env  env
	}

	setEnv := func(env env) {
		os.Setenv("DOMAIN", env.Domain)
		os.Setenv("ZEUS_TOKEN", env.Token)
		os.Setenv("ZEUS_API_URL", env.ApiUrl)
		os.Setenv("ZEUS_ADDITIONAL_RATE", env.AdditionalRate)
		os.Setenv("LOCALE_SPREADSHEET_ID", env.SpreadSheetId)
		os.Setenv("LOCALE_KEY_PREFIX", env.SheetKeyPrefix)
		os.Setenv("LOCALE_CREDENTIALS_FILE", env.Credentials)
		os.Setenv("LOCALE_SHEET_NAME", env.SheetName)
		os.Setenv("MAIL_HOST", env.Host)
		os.Setenv("MAIL_PORT", strconv.Itoa(env.Port))
		os.Setenv("MAIL_USERNAME", env.Username)
		os.Setenv("MAIL_PASSWORD", env.Password)
		os.Setenv("MAIL_FROM", env.From)

		os.Setenv("MAIL_TO", env.To)
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "test config",
			args: args{
				env: env{
					Domain:  "example.com",

					Token:  "token",
					ApiUrl: "url",
					AdditionalRate: "0.5",

					SpreadSheetId: "test",
					Credentials: "test",
					SheetName: "test",
					SheetKeyPrefix: "prefix",

					Host: "host",
					Port: 123,
					Username: "host",
					Password: "host",
					From: "host",

					To: "support@educate.market,poart1405@gmail.com",
				},
			},
			want: &Config{
				ZeusConfig{
					Token:  "token",
					ApiUrl: "url",
					AdditionalRate: 0.5,
				},
				LocaleConfig{
					SpreadSheetId: "test",
					Credentials: "test",
					SheetName: "test",
					SheetKeyPrefix: "prefix",
				},
				MailConfig{
					Host: "host",
					Port: 123,
					Username: "host",
					Password: "host",
					From: "host",

					To: []string{"support@educate.market", "poart1405@gmail.com"},
				},

				"example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.args.env)

			got, err := Init("../../.env")
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() Config error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() Config got = %v, want %v", got, tt.want)
			}
		})
	}
}
