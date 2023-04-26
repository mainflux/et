package api

import (
	"github.com/mainflux/callhome/internal/homing"
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

type saveTelemetryReq struct {
	homing.Telemetry
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
	offset    uint64
	limit     uint64
	repo      string
	IpAddress string `json:"ip_address"`
}

func (req listTelemetryReq) validate() error {
	if req.limit > maxLimitSize || req.limit < 1 {
		return ErrLimitSize
	}

	return nil
}
