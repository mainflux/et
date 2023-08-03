package api

import (
	"net/http"

	"github.com/mainflux/callhome"
	"github.com/mainflux/mainflux"
)

var _ mainflux.Response = (*uiRes)(nil)
var _ mainflux.Response = (*saveTelemetryRes)(nil)
var _ mainflux.Response = (*telemetryPageRes)(nil)
var _ mainflux.Response = (*telemetrySummaryRes)(nil)

type saveTelemetryRes struct {
	created bool
}

func (res saveTelemetryRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res saveTelemetryRes) Headers() map[string]string {
	if res.created {
		return map[string]string{}
	}

	return map[string]string{}
}

func (res saveTelemetryRes) Empty() bool {
	return true
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type telemetryPageRes struct {
	pageRes
	Telemetry []callhome.Telemetry `json:"telemetry"`
}

func (res telemetryPageRes) Code() int {
	return http.StatusOK
}

func (res telemetryPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res telemetryPageRes) Empty() bool {
	return false
}

type uiRes struct {
	code    int
	headers map[string]string
	html    []byte
}

// Code implements mainflux.Response
func (res uiRes) Code() int {
	if res.code == 0 {
		return http.StatusCreated
	}

	return res.code
}

// Empty implements mainflux.Response
func (res uiRes) Empty() bool {
	return res.html == nil
}

// Headers implements mainflux.Response
func (res uiRes) Headers() map[string]string {
	if res.headers == nil {
		return map[string]string{}
	}
	return res.headers
}

type telemetrySummaryRes struct {
	Countries        []callhome.CountrySummary `json:"countries,omitempty"`
	Cities           []callhome.CitySummary    `json:"cities,omitempty"`
	TotalDeployments int                       `json:"total_deployments,omitempty"`
}

// Code implements mainflux.Response
func (*telemetrySummaryRes) Code() int {
	return http.StatusOK
}

// Empty implements mainflux.Response
func (*telemetrySummaryRes) Empty() bool {
	return false
}

// Headers implements mainflux.Response
func (res *telemetrySummaryRes) Headers() map[string]string {
	return map[string]string{}
}
