package main

import (
	"context"
	"fmt"
	"log"
	"os"

	jaegerClient "github.com/mainflux/et/internal/clients/jaeger"
	"github.com/mainflux/et/internal/env"
	"github.com/mainflux/et/internal/homing"
	"github.com/mainflux/et/internal/homing/api"
	"github.com/mainflux/et/internal/homing/sheets"
	"github.com/mainflux/et/internal/server"
	httpserver "github.com/mainflux/et/internal/server/http"
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
	LogLevel       string `env:"MF_USERS_LOG_LEVEL"             envDefault:"info"`
	JaegerURL      string `env:"MF_JAEGER_URL"                  envDefault:"localhost:6831"`
	GCPCredFile    string `env:"MF_GCP_CRED"					envDefault:"sammydrive-7c852a28ee7f.json"`
	SpreadsheetId  string `env:"MF_SPREADSHEET_ID" 				envDefault:"1neq9yx6kEKx6HFWJepqhBBs_2qW0t56eE6Q5lZb-Hk8"`
	SheetId        int    `env:"MF_SHEET_ID"					envDefault:1`
	IPDatabaseFile string `env:"MF_IP_DB"						envDefault:"IP2LOCATION-LITE-DB5.BIN"`
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
		logger.Debug(fmt.Sprintf("failed to init logger: %s", err.Error()))
		return
	}

	tracer, closer, err := jaegerClient.NewTracer("users", cfg.JaegerURL)
	if err != nil {
		logger.Debug(fmt.Sprintf("failed to init Jaeger: %s", err))
		return
	}
	defer closer.Close()

	svc := newService(logger, cfg.IPDatabaseFile, cfg.GCPCredFile, cfg.SpreadsheetId, cfg.SheetId)

	httpServerConfig := server.Config{Port: defSvcHttpPort}
	if err := env.Parse(&httpServerConfig, env.Options{Prefix: envPrefixHttp, AltPrefix: envPrefix}); err != nil {
		logger.Debug(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err))
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

func newService(logger mflog.Logger, ipDB, credFile, spreadsheetID string, sheetID int) homing.Service {
	repo, err := sheets.New(credFile, spreadsheetID, sheetID)
	if err != nil {
		logger.Debug(err.Error())
		return nil
	}
	locsvc, err := homing.NewLocationService(ipDB)
	if err != nil {
		logger.Debug(err.Error())
		return nil
	}
	return homing.New(repo, locsvc)
}
