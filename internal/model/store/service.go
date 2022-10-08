package store

import (
	"sync"
	"time"
)

func New() *Service {
	return &Service{
		Categories: []category{},
		Purchases:  []purchase{},
	}
}

type Service struct {
	CategoriesMutex sync.Mutex
	PurchasesMutex  sync.Mutex
	UsersMutex      sync.Mutex

	Categories []category
	Purchases  []purchase
	Users      []user
}

type user struct {
	UserID   int64
	Currency Currency // выбранная пользователем валюта
}

type category struct {
	UserID   int64
	Category string
}

type purchase struct {
	UserID   int64
	Sum      float64 // сумма траты в рублях
	Category string
	Date     time.Time

	// коэффициенты валют на момент совершения траты
	USDRatio float64
	CNYRatio float64
	EURRatio float64
}

// Currency тип валюты
type Currency byte

const (
	// RUB валюта - рубль
	RUB Currency = 0

	// USD валюта - доллар
	USD Currency = 1

	// EUR валюта - евро
	EUR Currency = 2

	// CNY валюта - китайский юань
	CNY Currency = 3
)
