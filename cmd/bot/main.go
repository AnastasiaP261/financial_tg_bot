package main

import (
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/chart_drawing"
	"log"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/store"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	db := store.New()

	chartDrawingModel := chart_drawing.New()
	purchasesModel := purchases.New(db, chartDrawingModel)
	msgModel := messages.New(tgClient, purchasesModel)

	tgClient.ListenUpdates(msgModel)
}
