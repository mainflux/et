package callhome

import (
	"context"
	"time"

	"github.com/lib/pq"
)

type Telemetry struct {
	Services    pq.StringArray `json:"services,omitempty" db:"services"`
	Service     string         `json:"service,omitempty" db:"service"`
	Longitude   float64        `json:"longitude,omitempty" db:"longitude"`
	Latitude    float64        `json:"latitude,omitempty" db:"latitude"`
	IpAddress   string         `json:"-" db:"ip_address"`
	Version     string         `json:"mainflux_version,omitempty" db:"mf_version"`
	LastSeen    time.Time      `json:"last_seen" db:"service_time"`
	Country     string         `json:"country,omitempty" db:"country"`
	City        string         `json:"city,omitempty" db:"city"`
	ServiceTime time.Time      `json:"timestamp" db:"time"`
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

type TelemetrySummary struct {
	Countries   []string `json:"countries,omitempty"`
	IpAddresses []string `json:"ip_addresses,omitempty"`
}

// TelemetryRepository specifies an account persistence API.
type TelemetryRepo interface {
	// Save persists the telemetry event. A non-nil error is returned to indicate
	// operation failure.
	Save(ctx context.Context, t Telemetry) error

	// RetrieveAll retrieves all telemetry events.
	RetrieveAll(ctx context.Context, pm PageMetadata) (TelemetryPage, error)
	// RetrieveDistinctIPsCOuntries gets distinct ip addresses and countries from database.
	RetrieveDistinctIPsCountries(ctx context.Context) (TelemetrySummary, error)
}
