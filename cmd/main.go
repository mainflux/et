package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mainflux/callhome/callhome"
	"github.com/mainflux/callhome/callhome/api"
	"github.com/mainflux/callhome/callhome/repository/timescale"
	"github.com/mainflux/callhome/internal"
	jaegerClient "github.com/mainflux/callhome/internal/clients/jaeger"
	"github.com/mainflux/callhome/internal/clients/postgres"
	"github.com/mainflux/callhome/internal/env"
	"github.com/mainflux/callhome/internal/server"
	httpserver "github.com/mainflux/callhome/internal/server/http"
	mflog "github.com/mainflux/mainflux/logger"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "callhome"
	envPrefix      = "MF_CALLHOME_"
	envPrefixHttp  = "MF_CALLHOME_"
	defSvcHttpPort = "8855"
)

type config struct {
	LogLevel       string `env:"MF_CALLHOME_LOG_LEVEL"       envDefault:"info"`
	JaegerURL      string `env:"MF_JAEGER_URL"               envDefault:"localhost:6831"`
	IPDatabaseFile string `env:"MF_CALLHOME_IP_DB"           envDefault:"./IP2LOCATION-LITE-DB5.BIN"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err)
	}

	logger, err := mflog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to init logger: %s", err.Error()))
	}

	timescaleDB, err := postgres.Setup(envPrefix, timescale.Migration())
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to setup timescale db : %s", err))
	}

	tracer, closer, err := jaegerClient.NewTracer("users", cfg.JaegerURL)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to init Jaeger: %s", err))
	}
	defer closer.Close()

	svc, err := newService(ctx, logger, cfg.IPDatabaseFile, timescaleDB)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initialize service: %s", err.Error()))
		return
	}

	httpServerConfig := server.Config{Port: defSvcHttpPort}
	if err := env.Parse(&httpServerConfig, env.Options{Prefix: envPrefixHttp, AltPrefix: envPrefix}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err.Error()))
		return
	}
	hs := httpserver.New(ctx, cancel, svcName, httpServerConfig, api.MakeHandler(svc, tracer, logger), logger)

	g.Go(func() error {
		return hs.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, hs)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("%s service terminated: %s", svcName, err))
	}
}

func newService(ctx context.Context, logger mflog.Logger, ipDB string, timescaleDB *sqlx.DB) (callhome.Service, error) {
	timescaleRepo := timescale.New(timescaleDB)
	locSvc, err := callhome.NewLocationService(ipDB)
	if err != nil {
		return nil, err
	}
	svc := callhome.New(timescaleRepo, locSvc)
	counter, latency := internal.MakeMetrics(svcName, "api")
	svc = api.MetricsMiddleware(svc, counter, latency)
	svc = api.LoggingMiddleware(svc, logger)
	return svc, nil
}
