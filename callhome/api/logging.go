package api

import (
	"context"
	"fmt"
	"time"

	"github.com/mainflux/callhome/callhome"
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

// GetAll adds logging middleware to get all service.
func (lm *loggingMiddleware) Retrieve(ctx context.Context, repo string, pm callhome.PageMetadata) (telemetryPage callhome.TelemetryPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get all telemetry with took %s to complete", time.Since(begin))
		if err != nil {
			lm.hommingLogger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.hommingLogger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Retrieve(ctx, repo, pm)
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
