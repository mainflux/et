package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/callhome"
)

func saveEndpoint(svc callhome.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(saveTelemetryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		tel := callhome.Telemetry{
			Service:     req.Service,
			IpAddress:   req.IpAddress,
			Version:     req.Version,
			ServiceTime: req.LastSeen,
		}
		if err := svc.Save(ctx, tel); err != nil {
			return nil, err
		}
		res := saveTelemetryRes{
			created: true,
		}
		return res, nil
	}
}

func retrieveEndpoint(svc callhome.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listTelemetryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		pm := callhome.PageMetadata{
			Offset: req.offset,
			Limit:  req.limit,
		}
		filter := callhome.TelemetryFilters{
			From:    req.from,
			To:      req.to,
			Country: req.country,
			City:    req.city,
			Version: req.version,
		}
		tm, err := svc.Retrieve(ctx, pm, filter)
		if err != nil {
			return nil, err
		}
		res := telemetryPageRes{
			pageRes: pageRes{
				Total:  tm.Total,
				Offset: tm.Offset,
				Limit:  tm.Limit,
			},
			Telemetry: tm.Telemetry,
		}
		return res, nil
	}
}

func retrieveSummaryEndpoint(svc callhome.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listTelemetryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		filter := callhome.TelemetryFilters{
			From:    req.from,
			To:      req.to,
			Country: req.country,
			City:    req.city,
			Version: req.version,
		}
		summary, err := svc.RetrieveSummary(ctx, filter)
		if err != nil {
			return nil, err
		}
		return telemetrySummaryRes{
			Countries:        summary.Countries,
			TotalDeployments: summary.TotalDeployments,
		}, nil
	}
}

func serveUI(svc callhome.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listTelemetryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		filter := callhome.TelemetryFilters{
			From:    req.from,
			To:      req.to,
			Country: req.country,
			City:    req.city,
			Version: req.version,
		}
		res, err := svc.ServeUI(ctx, filter)
		return uiRes{
			html: res,
		}, err
	}
}
