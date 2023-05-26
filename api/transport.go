package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/callhome"
	"github.com/mainflux/callhome/timescale"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/pkg/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType  = "application/json"
	offsetKey    = "offset"
	limitKey     = "limit"
	ipAddressKey = "ip"
	metadataKey  = "metadata"
	statusKey    = "status"
	defOffset    = 0
	defLimit     = 10
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc callhome.Service, tracer opentracing.Tracer, logger logger.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(LoggingErrorEncoder(logger, encodeError)),
	}

	mux := bone.New()

	mux.Post("/telemetry", kithttp.NewServer(
		kitot.TraceServer(tracer, "save")(saveEndpoint(svc)),
		decodeSaveTelemetryReq,
		encodeResponse,
		opts...,
	))

	mux.Get("/telemetry", kithttp.NewServer(
		kitot.TraceServer(tracer, "retrieve")(retrieveEndpoint(svc)),
		decodeRetrieve,
		encodeResponse,
		opts...,
	))

	mux.Get("/telemetry/summary", kithttp.NewServer(
		kitot.TraceServer(tracer, "retrieve-summary")(retrieveSummaryEndpoint(svc)),
		kithttp.NopRequestDecoder,
		encodeResponse,
		opts...,
	))
	mux.GetFunc("/health", mainflux.Health("telemetry"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch {
	case errors.Contains(err, errors.ErrInvalidQueryParams),
		errors.Contains(err, errors.ErrMalformedEntity),
		err == ErrLimitSize,
		err == ErrOffsetSize:
		w.WriteHeader(http.StatusBadRequest)
	case errors.Contains(err, timescale.ErrInvalidEvent):
		w.WriteHeader(http.StatusForbidden)
	case errors.Contains(err, errors.ErrUnsupportedContentType):
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errors.Contains(err, uuid.ErrGeneratingID):
		w.WriteHeader(http.StatusInternalServerError)
	case errors.Contains(err, timescale.ErrSaveEvent),
		errors.Contains(err, timescale.ErrTransRollback):
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if errorVal, ok := err.(errors.Error); ok {
		w.Header().Set("Content-Type", contentType)
		if err := json.NewEncoder(w).Encode(ErrorRes{Err: errorVal.Msg()}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func decodeRetrieve(_ context.Context, r *http.Request) (interface{}, error) {
	o, err := ReadUintQuery(r, offsetKey, defOffset)
	if err != nil {
		return nil, err
	}

	l, err := ReadUintQuery(r, limitKey, defLimit)
	if err != nil {
		return nil, err
	}

	req := listTelemetryReq{
		offset: o,
		limit:  l,
	}
	return req, nil
}

func decodeSaveTelemetryReq(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}

	var telemetry saveTelemetryReq
	if err := json.NewDecoder(r.Body).Decode(&telemetry); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}

	return telemetry, nil
}
