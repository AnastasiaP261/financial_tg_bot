package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/metrics"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/fixer"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/redis"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	tgmsghandler "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg/messages"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/env"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/chart_drawing"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/db"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/exchange_rates"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	logswrapper "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/logs"
)

const (
	port = 8080
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
	db, err := db.New(ctx, config)
	if err != nil {
		log.Fatal("database init failed:", err)
	}

	redis, err := redis.New(ctx, config)
	if err != nil {
		log.Fatal("redis init failed:", err)
	}

	// CLIENTS
	msgHandler := initTgMsgHandler(config)
	fixerClient := fixer.New(ctx, config)

	// MODELS
	chartDrawingModel := chart_drawing.New()
	exchangesRatesModel := exchange_rates.New(fixerClient)
	purchasesModel := purchases.New(db, chartDrawingModel, exchangesRatesModel)

	msgModel := messages.New(msgHandler, purchasesModel, redis)

	// INFRA
	logs.InitLogger()

	// ПОЕХАЛИ!!
	errG, ctx := errgroup.WithContext(ctx)
	errG.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatal("error starting http server:", err)
			return err
		}
		logs.Info("starting http server", zap.Int("port", port))

		return nil
	})

	errG.Go(func() error {
		listener, err := tg.New(config, msgHandler)
		if err != nil {
			log.Fatal("tg listener init failed")
			return err
		}
		logs.Info("messages listen started")
		listener.ListenUpdates(ctx, msgModel)
		return nil
	})

	if err = errG.Wait(); err != nil {
		log.Fatal("errgroup error:", err)
	}
}

func initTgMsgHandler(conf *config.Service) *metrics.Wrapper {
	msgHandler, err := tgmsghandler.New(conf)
	if err != nil {
		log.Fatal("tg msg handler init failed:", err)
	}
	return metrics.NewWrapper(logswrapper.NewWrapper(msgHandler))
}
