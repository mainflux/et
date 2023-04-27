package sheets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mainflux/callhome/callhome"
	"github.com/mainflux/callhome/callhome/repository"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ callhome.TelemetryRepo = (*repo)(nil)

const (
	sheetRange = "Sheet1!A:H"
	sheetsAuth = "https://www.googleapis.com/auth/spreadsheets"
)

var (
	errFailedStringConversion error  = errors.New("failed to convert field to string")
	errFailedFloatConversion  error  = errors.New("failed to convert field to float64")
	errFailedParsingLastSeen  error  = errors.New("failed to parse last seen")
	errInvalidRowLength       string = "invalid row length: expected 8, got %d"
)

type repo struct {
	sheetsSvc     *sheets.Service
	sheetName     string
	spreadsheetId string
}

// New Creates a new telemetry repo using google sheets.
func New(credFile, spreadsheetId string, sheetID int) (callhome.TelemetryRepo, error) {
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

// RetrieveAll gets all records from repo.
func (r repo) RetrieveAll(ctx context.Context, pm callhome.PageMetadata) (callhome.TelemetryPage, error) {
	resp, err := r.sheetsSvc.Spreadsheets.Values.Get(r.spreadsheetId, sheetRange).Do()
	if err != nil {
		return callhome.TelemetryPage{}, err
	}
	var telPage callhome.TelemetryPage
	telPage.PageMetadata = pm
	for _, row := range resp.Values {
		tel, err := fromRow(row)
		if err != nil {
			return telPage, err
		}
		telPage.Telemetry = append(telPage.Telemetry, tel)
	}
	return telPage, nil
}

// RetrieveByIP get record by ip address.
func (r *repo) RetrieveByIP(ctx context.Context, ip string) (callhome.Telemetry, error) {
	resp, err := r.sheetsSvc.Spreadsheets.Values.Get(r.spreadsheetId, sheetRange).Do()
	if err != nil {
		return callhome.Telemetry{}, err
	}
	for _, row := range resp.Values {
		if len(row) >= 2 && row[1] == ip {
			tel, err := fromRow(row)
			return tel, err
		}
	}
	return callhome.Telemetry{}, repository.ErrRecordNotFound
}

// Save adds record to repo.
func (r repo) Save(ctx context.Context, t callhome.Telemetry) error {
	rrow, err := toRow(t)
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

// Update update record to repo.
func (r repo) Update(ctx context.Context, t callhome.Telemetry) error {
	rrow, err := toRow(t)
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

// ToRow converts telemetry event to google sheets row.
func toRow(t callhome.Telemetry) ([]interface{}, error) {
	lastSeen, err := t.LastSeen.MarshalText()
	if err != nil {
		return nil, err
	}
	return []interface{}{t.IpAddress, t.Latitude, t.Longitude, strings.Join(t.Services, ","), t.Version, string(lastSeen), t.City, t.Country}, nil
}

// FromRow converts a Google Sheets row to a Telemetry struct.
func fromRow(row []interface{}) (callhome.Telemetry, error) {
	var t callhome.Telemetry
	if len(row) != 8 {
		return callhome.Telemetry{}, fmt.Errorf(errInvalidRowLength, len(row))
	}
	ipAddress, ok := row[0].(string)
	if !ok {
		return callhome.Telemetry{}, errFailedStringConversion
	}
	t.IpAddress = ipAddress
	lat, err := strconv.ParseFloat(row[1].(string), 64)
	if err != nil {
		return callhome.Telemetry{}, errors.Join(errFailedFloatConversion, err)
	}
	t.Latitude = lat
	long, err := strconv.ParseFloat(row[2].(string), 64)
	if err != nil {
		return callhome.Telemetry{}, errors.Join(errFailedFloatConversion, err)
	}
	t.Longitude = long
	services, ok := row[3].(string)
	if !ok {
		return callhome.Telemetry{}, errFailedStringConversion
	}
	t.Services = strings.Split(services, ",")
	version, ok := row[4].(string)
	if !ok {
		return callhome.Telemetry{}, errFailedStringConversion
	}
	t.Version = version
	lastSeen, ok := row[5].(string)
	if !ok {
		return callhome.Telemetry{}, errFailedStringConversion
	}
	if err = t.LastSeen.UnmarshalText([]byte(lastSeen)); err != nil {
		return callhome.Telemetry{}, errors.Join(errFailedParsingLastSeen, err)
	}
	city, ok := row[6].(string)
	if !ok {
		return callhome.Telemetry{}, errFailedStringConversion
	}
	t.City = city
	country, ok := row[7].(string)
	if !ok {
		return callhome.Telemetry{}, errFailedStringConversion
	}
	t.Country = country
	return t, nil
}
