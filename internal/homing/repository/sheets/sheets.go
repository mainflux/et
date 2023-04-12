package sheets

import (
	"context"
	"os"

	"github.com/mainflux/et/internal/homing"
	"github.com/mainflux/et/internal/homing/repository"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ homing.TelemetryRepo = (*repo)(nil)

const (
	sheetRange = "Sheet1!A:I"
	sheetsAuth = "https://www.googleapis.com/auth/spreadsheets"
)

type repo struct {
	sheetsSvc     *sheets.Service
	sheetName     string
	spreadsheetId string
}

// New Creates a new telemetry repo using google sheets.
func New(credFile, spreadsheetId string, sheetID int) (homing.TelemetryRepo, error) {
	credBytes, err := os.ReadFile(credFile)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(credBytes, sheetsAuth)
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background())

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	res, err := srv.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return nil, err
	}

	sheetName := ""
	for _, v := range res.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetID) {
			sheetName = prop.Title
			break
		}
	}

	return &repo{
		sheetsSvc:     srv,
		sheetName:     sheetName,
		spreadsheetId: spreadsheetId,
	}, nil
}

// RetrieveAll implements homing.TelemetryRepo.
func (r repo) RetrieveAll(ctx context.Context, pm homing.PageMetadata) (homing.TelemetryPage, error) {

	resp, err := r.sheetsSvc.Spreadsheets.Values.Get(r.spreadsheetId, sheetRange).Do()
	if err != nil {
		return homing.TelemetryPage{}, err
	}
	var telPage homing.TelemetryPage
	telPage.PageMetadata = pm
	for _, row := range resp.Values {
		var tel homing.Telemetry
		if err := tel.FromRow(row); err != nil {
			return telPage, err
		}
		telPage.Telemetry = append(telPage.Telemetry, tel)
	}
	return telPage, nil
}

// RetrieveByIP implements homing.TelemetryRepo.
func (r *repo) RetrieveByIP(ctx context.Context, ip string) (homing.Telemetry, error) {
	resp, err := r.sheetsSvc.Spreadsheets.Values.Get(r.spreadsheetId, sheetRange).Do()
	if err != nil {
		return homing.Telemetry{}, err
	}
	for _, row := range resp.Values {
		if len(row) >= 2 && row[1] == ip {
			var tel homing.Telemetry
			err = tel.FromRow(row)
			return tel, err
		}
	}
	return homing.Telemetry{}, repository.ErrRecordNotFound
}

// Save implements homing.TelemetryRepo.
func (r repo) Save(ctx context.Context, t homing.Telemetry) error {
	rrow, err := t.ToRow()
	if err != nil {
		return err
	}
	row := &sheets.ValueRange{
		Values: [][]interface{}{rrow},
	}

	res, err := r.sheetsSvc.Spreadsheets.Values.Append(r.spreadsheetId, r.sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		return err
	}
	return nil
}

// UpdateTelemetry implements homing.TelemetryRepo.
func (r repo) UpdateTelemetry(ctx context.Context, t homing.Telemetry) error {
	rrow, err := t.ToRow()
	if err != nil {
		return err
	}
	updateValueRange := &sheets.ValueRange{
		Values: [][]interface{}{rrow},
	}
	if _, err := r.sheetsSvc.Spreadsheets.Values.Update(r.spreadsheetId, sheetRange, updateValueRange).ValueInputOption("USER_ENTERED").Do(); err != nil {
		return err
	}
	return nil
}
