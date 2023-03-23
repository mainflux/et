package main

import (
	"context"
	"fmt"
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
	JaegerURL      = "localhost:6831"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	logger, err := mflog.New(os.Stdout, "info")
	if err != nil {
		logger.Debug(fmt.Sprintf("failed to init logger: %s", err.Error()))
	}

	tracer, closer, err := jaegerClient.NewTracer("users", JaegerURL)
	if err != nil {
		logger.Debug(fmt.Sprintf("failed to init Jaeger: %s", err))
	}
	defer closer.Close()

	svc := newService(logger)

	httpServerConfig := server.Config{Port: defSvcHttpPort}
	if err := env.Parse(&httpServerConfig, env.Options{Prefix: envPrefixHttp, AltPrefix: envPrefix}); err != nil {
		logger.Debug(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err))
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

func newService(logger mflog.Logger) homing.Service {
	repo, err := sheets.New("", "", "")
	if err != nil {
		logger.Debug(err.Error())
	}
	locsvc, err := homing.NewLocationService("")
	if err != nil {
		logger.Debug(err.Error())
	}
	return homing.New(repo, locsvc)
}
