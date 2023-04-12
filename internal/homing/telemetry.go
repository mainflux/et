package homing

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Telemetry struct {
	ID        string    `json:"id,omitempty" db:"id"`
	Services  []string  `json:"services,omitempty" db:"-"`
	Service   string    `json:"service" db:"service"`
	Longitude float64   `json:"longitude,omitempty" db:"longitude"`
	Latitude  float64   `json:"latitude,omitempty" db:"latitude"`
	IpAddress string    `json:"ip_address" db:"ip_address"`
	Version   string    `json:"mainflux_version,omitempty" db:"mf_version"`
	LastSeen  time.Time `json:"last_seen" db:"last_seen"`
	Country   string    `json:"country,omitempty" db:"country"`
	City      string    `json:"city,omitempty" db:"city"`
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

	// Update updates Telemetry event.
	UpdateTelemetry(ctx context.Context, u Telemetry) error

	// RetrieveByIP retrieves telemetry by its unique identifier (i.e. ip address).
	RetrieveByIP(ctx context.Context, email string) (Telemetry, error)

	// RetrieveAll retrieves all telemetry events.
	RetrieveAll(ctx context.Context, pm PageMetadata) (TelemetryPage, error)
}

// ToRow converts telemetry event to google sheets row.
func (t *Telemetry) ToRow() ([]interface{}, error) {
	lastSeen, err := t.LastSeen.MarshalText()
	if err != nil {
		return nil, err
	}
	return []interface{}{t.ID, t.IpAddress, t.Latitude, t.Longitude, strings.Join(t.Services, ","), t.Version, string(lastSeen), t.City, t.Country}, nil
}

// FromRow converts a Google Sheets row to a Telemetry struct.
func (t *Telemetry) FromRow(row []interface{}) error {
	if len(row) != 9 {
		return fmt.Errorf("invalid row length: expected 6, got %d", len(row))
	}
	id, ok := row[0].(string)
	if !ok {
		return errors.New("failed to convert ID to string")
	}
	t.ID = id
	ipAddress, ok := row[1].(string)
	if !ok {
		return errors.New("failed to convert IP address to string")
	}
	t.IpAddress = ipAddress
	lat, err := strconv.ParseFloat(row[2].(string), 64)
	if err != nil {
		return fmt.Errorf("failed to convert latitude to float64: %v", err)
	}
	t.Latitude = lat
	long, err := strconv.ParseFloat(row[3].(string), 64)
	if err != nil {
		return fmt.Errorf("failed to convert longitude to float64: %v", err)
	}
	t.Longitude = long
	services, ok := row[4].(string)
	if !ok {
		return errors.New("failed to convert services to string")
	}
	t.Services = strings.Split(services, ",")
	version, ok := row[5].(string)
	if !ok {
		return errors.New("failed to convert version to string")
	}
	t.Version = version
	lastSeen, ok := row[6].(string)
	if !ok {
		return errors.New("failed to convert lastSeen to string")
	}
	if err = t.LastSeen.UnmarshalText([]byte(lastSeen)); err != nil {
		return fmt.Errorf("failed to parse last seen: %v", err)
	}
	city, ok := row[7].(string)
	if !ok {
		return errors.New("failed to convert ID to string")
	}
	t.City = city
	country, ok := row[8].(string)
	if !ok {
		return errors.New("failed to convert ID to string")
	}
	t.Country = country
	return nil
}
