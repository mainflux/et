package sheets

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/mainflux/et/internal/homing"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ homing.TelemetryRepo = (*repo)(nil)

type repo struct {
	sheetsSvc     *sheets.Service
	sheetName     string
	spreadsheetId string
}

func New(tokenFile, sheetName, spreadsheetId string) (homing.TelemetryRepo, error) {
	//credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	if err != nil {
		return nil, err
	}
	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	// create client with config and context
	client := config.Client(context.Background())

	// create new service using client
	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &repo{
		sheetsSvc:     srv,
		sheetName:     sheetName,
		spreadsheetId: spreadsheetId,
	}, nil
}

// RetrieveAll implements homing.TelemetryRepo
func (*repo) RetrieveAll(ctx context.Context, pm homing.PageMetadata) ([]homing.Telemetry, error) {
	return nil, nil
}

// RetrieveByIP implements homing.TelemetryRepo
func (r *repo) RetrieveByIP(ctx context.Context, email string) (*homing.Telemetry, error) {
	//res, err := r.sheetsSvc.Spreadsheets.Values.Get(r.spreadsheetId, "B:B").Do()
	return nil, nil

}

// Save implements homing.TelemetryRepo
func (r *repo) Save(ctx context.Context, t homing.Telemetry) error {
	t.ID = uuid.New().String()
	// Append value to the sheet.
	row := &sheets.ValueRange{
		Values: [][]interface{}{{t.ID, t.IpAddress, t.Latitutde, t.Longitude, strings.Join(t.Services, ","), t.Version}},
	}

	response2, err := r.sheetsSvc.Spreadsheets.Values.Append(r.spreadsheetId, r.sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || response2.HTTPStatusCode != 200 {
		return err
	}
	return nil
}

// UpdateTelemetry implements homing.TelemetryRepo
func (*repo) UpdateTelemetry(ctx context.Context, u homing.Telemetry) error {
	return nil
}
