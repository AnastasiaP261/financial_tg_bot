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

	Categories []category
	Purchases  []purchase
}

type category struct {
	User     string
	Category string
}

type purchase struct {
	User     string
	Sum      float64
	Category string
	Date     time.Time
}
