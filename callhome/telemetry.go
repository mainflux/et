package callhome

import (
	"context"
	"time"
)

type Telemetry struct {
	Services    []string  `json:"services,omitempty" db:"-"`
	Service     string    `json:"service,omitempty" db:"service"`
	Longitude   float64   `json:"longitude,omitempty" db:"longitude"`
	Latitude    float64   `json:"latitude,omitempty" db:"latitude"`
	IpAddress   string    `json:"ip_address" db:"ip_address"`
	Version     string    `json:"mainflux_version,omitempty" db:"mf_version"`
	LastSeen    time.Time `json:"last_seen" db:"time"`
	Country     string    `json:"country,omitempty" db:"country"`
	City        string    `json:"city,omitempty" db:"city"`
	ServiceTime time.Time `json:"service_last_seen" db:"service_time"`
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
	Update(ctx context.Context, u Telemetry) error

	// RetrieveByIP retrieves telemetry by its unique identifier (i.e. ip address).
	RetrieveByIP(ctx context.Context, email string) (Telemetry, error)

	// RetrieveAll retrieves all telemetry events.
	RetrieveAll(ctx context.Context, pm PageMetadata) (TelemetryPage, error)
}
