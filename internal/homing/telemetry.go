package homing

import (
	"context"
	"strconv"
	"strings"
)

type Telemetry struct {
	ID        string   `json:"-"`
	Services  []string `json:"-"`
	Longitude float64  `json:"-"`
	Latitutde float64  `json:"-"`
	IpAddress string   `json:"ip_address"`
	Version   string   `json:"mainflux_version"`
}

type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

type TelemetryPage struct {
	PageMetadata
	Telemetry []Telemetry
}

// TelemetryRepository specifies an account persistence API.
type TelemetryRepo interface {
	// Save persists the telemetry event. A non-nil error is returned to indicate
	// operation failure.
	Save(ctx context.Context, t Telemetry) error

	// Update updates Telemetry event
	UpdateTelemetry(ctx context.Context, u Telemetry, row int) error

	// RetrieveByIP retrieves telemetry by its unique identifier (i.e. ip address).
	RetrieveByIP(ctx context.Context, email string) (*Telemetry, int, error)

	// RetrieveAll retrieves all telemetry for given array of telemetry IDs.
	RetrieveAll(ctx context.Context, pm PageMetadata) ([]Telemetry, error)
}

func (t *Telemetry) ToRow() []interface{} {
	return []interface{}{t.ID, t.IpAddress, t.Latitutde, t.Longitude, strings.Join(t.Services, ","), t.Version}
}

func (t *Telemetry) FromRow(row []interface{}) {
	t.ID = row[0].(string)
	t.IpAddress = row[1].(string)
	lat, _ := strconv.ParseFloat(row[2].(string), 64)
	t.Latitutde = lat
	long, _ := strconv.ParseFloat(row[3].(string), 64)
	t.Longitude = long
	t.Version = row[5].(string)
	t.Services = strings.Split(row[4].(string), ",")
}
