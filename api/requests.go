package api

import (
	"time"

	"github.com/mainflux/mainflux/pkg/errors"
)

var (
	// ErrLimitSize indicates that an invalid limit.
	ErrLimitSize = errors.New("invalid limit size")
	// ErrOffsetSize indicates an invalid offset.
	ErrOffsetSize = errors.New("invalid offset size")
	// ErrInvalidDateRange indicates date from and to are invalid.
	ErrInvalidDateRange = errors.New("invalid date range")
)

const maxLimitSize = 100

type saveTelemetryReq struct {
	Service   string    `json:"service"`
	IpAddress string    `json:"ip_address"`
	Version   string    `json:"mainflux_version"`
	LastSeen  time.Time `json:"last_seen"`
}

func (req saveTelemetryReq) validate() error {
	if req.Service == "" {
		return errors.ErrMalformedEntity
	}

	if req.IpAddress == "" {
		return errors.ErrMalformedEntity
	}
	if req.Version == "" {
		return errors.ErrMalformedEntity
	}

	return nil
}

type listTelemetryReq struct {
	offset  uint64
	limit   uint64
	from    time.Time
	to      time.Time
	country string
	city    string
}

func (req listTelemetryReq) validate() error {
	if req.limit > maxLimitSize || req.limit < 1 {
		return ErrLimitSize
	}

	if !req.from.IsZero() && req.to.Before(req.from) {
		return ErrInvalidDateRange
	}

	return nil
}
