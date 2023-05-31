package api

import (
	"context"
	"fmt"
	"time"

	"github.com/mainflux/callhome"
	"github.com/mainflux/mainflux/logger"
)

var _ callhome.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	hommingLogger logger.Logger
	svc           callhome.Service
}

// LoggingMiddleware is a middleware that adds logging facilities to the core homing service.
func LoggingMiddleware(svc callhome.Service, logger logger.Logger) callhome.Service {
	return &loggingMiddleware{logger, svc}
}

// Retrieve adds logging middleware to retrieve service.
func (lm *loggingMiddleware) Retrieve(ctx context.Context, pm callhome.PageMetadata) (telemetryPage callhome.TelemetryPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method retrieve with took %s to complete", time.Since(begin))
		if err != nil {
			lm.hommingLogger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.hommingLogger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Retrieve(ctx, pm)
}

// Save adds logging middleware to save service.
func (lm *loggingMiddleware) Save(ctx context.Context, t callhome.Telemetry) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method save telemetry event took %s to complete", time.Since(begin))
		if err != nil {
			lm.hommingLogger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.hommingLogger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Save(ctx, t)
}

func (lm *loggingMiddleware) RetrieveSummary(ctx context.Context) (summary callhome.TelemetrySummary, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method retrieve summary event took %s to complete", time.Since(begin))
		if err != nil {
			lm.hommingLogger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.hommingLogger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.RetrieveSummary(ctx)
}
