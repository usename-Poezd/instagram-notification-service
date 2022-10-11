package googlesheets

import (
	"context"
	"io/ioutil"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Sheet interface {
	Get(rang string) (*sheets.ValueRange, error)
}


type GoogleSheet struct {
	service *sheets.Service
	spreadsheetId string
}

func NewGoogleSheet(ctx context.Context, credentialsFilePath string, spreadsheetId string) Sheet {
	b, err := ioutil.ReadFile(credentialsFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	service, err := sheets.NewService(ctx, option.WithCredentialsJSON(b), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
			log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return &GoogleSheet{
		service,
		spreadsheetId,
	}
}


func (gs *GoogleSheet) Get(rang string) (*sheets.ValueRange, error) {
	return gs.service.Spreadsheets.Values.Get(gs.spreadsheetId, rang).Do()
}


