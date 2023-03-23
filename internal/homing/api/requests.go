package api

import (
	"github.com/mainflux/et/internal/homing"
	"github.com/mainflux/mainflux/pkg/errors"
)

var (
	// ErrBearerToken indicates missing or invalid bearer user token.
	ErrBearerToken = errors.New("missing or invalid bearer user token")
	// ErrLimitSize indicates that an invalid limit.
	ErrLimitSize = errors.New("invalid limit size")
	// ErrOffsetSize indicates an invalid offset.
	ErrOffsetSize = errors.New("invalid offset size")
)

const maxLimitSize = 100

type telemetryReq struct {
	homing.Telemetry
	ServiceName string `json:"service"`
}

func (req telemetryReq) validate() error {
	if req.ServiceName == "" {
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
	token     string
	offset    uint64
	limit     uint64
	IpAddress string `json:"ip_address"`
}

func (req listTelemetryReq) validate() error {
	if req.token == "" {
		return ErrBearerToken
	}

	if req.limit > maxLimitSize || req.limit < 1 {
		return ErrLimitSize
	}

	return nil
}
