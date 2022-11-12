package main

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/redis"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/kafka/consumer"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/env"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/db"
	tracing "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/tracing"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	port        = 8081
	serviceName = "financial-reports"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("unknown environment")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := env.SetEnvVariable(os.Args[1]); err != nil {
		log.Fatal("environment variable set failed:", err)
	}

	// CONFIG
	config, err := config.New(env.GetConfigFilePath())
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	// REPOSITORIES
	_, err = db.New(ctx, config)
	if err != nil {
		log.Fatal("database init failed:", err)
	}
	_, err = redis.New(ctx, config)
	if err != nil {
		log.Fatal("redis init failed:", err)
	}

	// MODELS

	// INFRA
	logger := logs.InitLogger()
	tracing.InitTracing(logger, serviceName)

	// ПОЕХАЛИ!!
	errG, ctx := errgroup.WithContext(ctx)
	errG.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())

		logs.Info("starting http server", zap.Int("port", port))
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatal("error starting http server:", err)
			return err
		}
		return nil
	})

	errG.Go(func() error {
		if err = consumer.StartConsumerGroup(ctx); err != nil {
			log.Fatal("starting consumer group error:", err)
			return err
		}
		return nil
	})

	if err = errG.Wait(); err != nil {
		log.Fatal("errgroup error:", err)
	}
}
