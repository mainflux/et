package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mainflux/et/internal"
	authClient "github.com/mainflux/et/internal/clients/grpc/auth"
	jaegerClient "github.com/mainflux/et/internal/clients/jaeger"
	"github.com/mainflux/et/internal/env"
	"github.com/mainflux/et/internal/homing"
	"github.com/mainflux/et/internal/homing/api"
	"github.com/mainflux/et/internal/homing/repository/sheets"
	"github.com/mainflux/et/internal/homing/repository/timescale"
	"github.com/mainflux/et/internal/server"
	httpserver "github.com/mainflux/et/internal/server/http"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "homing"
	envPrefix      = "MF_HOMING_"
	envPrefixHttp  = "MF_HOMING_"
	defSvcHttpPort = "8080"
)

type config struct {
	LogLevel       string `env:"MF_USERS_LOG_LEVEL"  envDefault:"info"`
	JaegerURL      string `env:"MF_JAEGER_URL"       envDefault:"localhost:6831"`
	GCPCredFile    string `env:"MF_GCP_CRED"`
	SpreadsheetId  string `env:"MF_SPREADSHEET_ID"`
	SheetId        int    `env:"MF_SHEET_ID"         envDefault:"0"`
	IPDatabaseFile string `env:"MF_IP_DB"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err)
	}

	timescaleConf := timescale.Config{}
	if err := env.Parse(&timescaleConf); err != nil {
		log.Fatalf("failed to load %s timescale configuration : %s", svcName, err)
	}

	timescaleDB, err := timescale.Connect(timescaleConf)
	if err != nil {
		log.Fatalf("failed to connect to timescale db : %s", err)
	}

	logger, err := mflog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to init logger: %s", err.Error()))
	}

	tracer, closer, err := jaegerClient.NewTracer("users", cfg.JaegerURL)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to init Jaeger: %s", err))
	}
	defer closer.Close()

	// Setup new auth grpc client
	auth, authHandler, err := authClient.Setup(envPrefix, cfg.JaegerURL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer authHandler.Close()
	logger.Info("Successfully connected to auth grpc server " + authHandler.Secure())

	svc, err := newService(logger, cfg.IPDatabaseFile, cfg.GCPCredFile, cfg.SpreadsheetId, cfg.SheetId, auth, timescaleDB)

	if err != nil {
		log.Printf("failed to initialize service: %s", err.Error())
		return
	}

	httpServerConfig := server.Config{Port: defSvcHttpPort}
	if err := env.Parse(&httpServerConfig, env.Options{Prefix: envPrefixHttp, AltPrefix: envPrefix}); err != nil {
		log.Printf("failed to load %s HTTP server configuration : %s", svcName, err.Error())
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
		logger.Error(fmt.Sprintf("HTTP adapter service terminated: %s", err))
	}
}

func newService(logger mflog.Logger, ipDB, credFile, spreadsheetID string, sheetID int, auth mainflux.AuthServiceClient, timescaleDB *sqlx.DB) (homing.Service, error) {
	repo, err := sheets.New(credFile, spreadsheetID, sheetID)
	if err != nil {
		return nil, err
	}
	timescaleRepo := timescale.New(timescaleDB)
	locSvc, err := homing.NewLocationService(ipDB)
	if err != nil {
		return nil, err
	}
	svc := homing.New(timescaleRepo, repo, locSvc, auth)
	counter, latency := internal.MakeMetrics(svcName, "api")
	svc = api.MetricsMiddleware(svc, counter, latency)
	svc = api.LoggingMiddleware(svc, logger)
	return svc, nil
}
