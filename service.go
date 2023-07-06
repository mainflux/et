package callhome

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"
	"time"
)

const pageLimit = 1000

// Service Service to receive homing telemetry data, persist and retrieve it.
type Service interface {
	// Save saves the homing telemetry data and its location information.
	Save(ctx context.Context, t Telemetry) error
	// Retrieve retrieves homing telemetry data from the specified repository.
	Retrieve(ctx context.Context, pm PageMetadata) (TelemetryPage, error)
	// RetrieveSummary gets distinct countries and ip addresses
	RetrieveSummary(ctx context.Context) (TelemetrySummary, error)
	// ServeUI gets the callhome index html page
	ServeUI(ctx context.Context) ([]byte, error)
}

var _ Service = (*telemetryService)(nil)

type telemetryService struct {
	repo   TelemetryRepo
	locSvc LocationService
}

// New creates a new instance of the telemetry service.
func New(repo TelemetryRepo, locSvc LocationService) Service {
	return &telemetryService{
		repo:   repo,
		locSvc: locSvc,
	}
}

// Retrieve retrieves homing telemetry data from the specified repository.
func (ts *telemetryService) Retrieve(ctx context.Context, pm PageMetadata) (TelemetryPage, error) {
	return ts.repo.RetrieveAll(ctx, pm)
}

// Save saves the homing telemetry data and its location information.
func (ts *telemetryService) Save(ctx context.Context, t Telemetry) error {
	locRec, err := ts.locSvc.GetLocation(ctx, t.IpAddress)
	if err != nil {
		return err
	}
	t.City = locRec.City
	t.Country = locRec.Country_long
	t.Latitude = float64(locRec.Latitude)
	t.Longitude = float64(locRec.Longitude)
	t.LastSeen = time.Now()
	return ts.repo.Save(ctx, t)
}

func (ts *telemetryService) RetrieveSummary(ctx context.Context) (TelemetrySummary, error) {
	return ts.repo.RetrieveDistinctIPsCountries(ctx)
}

// ServeUI gets the callhome index html page
func (ts *telemetryService) ServeUI(ctx context.Context) ([]byte, error) {
	tmpl := template.Must(template.ParseFiles("./web/template/index.html"))

	summary, err := ts.repo.RetrieveDistinctIPsCountries(ctx)
	if err != nil {
		return nil, err
	}
	telPage, err := ts.repo.RetrieveAll(ctx, PageMetadata{Limit: pageLimit})
	if err != nil {
		return nil, err
	}

	pg, err := json.Marshal(telPage)
	if err != nil {
		return nil, err
	}
	data := struct {
		Countries     []string
		NoDeployments int
		NoCountries   int
		MapData       string
	}{
		Countries:     summary.Countries,
		NoDeployments: len(summary.IpAddresses),
		NoCountries:   len(summary.Countries),
		MapData:       string(pg),
	}

	var res bytes.Buffer
	if err = tmpl.Execute(&res, data); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
