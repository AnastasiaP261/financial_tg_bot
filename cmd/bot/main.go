package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/fixer"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/chart_drawing"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/exchange_rates"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/store"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(conf)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}
	fixerClient := fixer.New(ctx, conf)

	db := store.New()

	chartDrawingModel := chart_drawing.New()
	exchangesRatesModel := exchange_rates.New(fixerClient)
	purchasesModel := purchases.New(db, chartDrawingModel, exchangesRatesModel)

	msgModel := messages.New(tgClient, purchasesModel)

	tgClient.ListenUpdates(msgModel)
}
