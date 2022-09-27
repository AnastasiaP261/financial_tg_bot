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
	UserID   int64
	Category string
}

type purchase struct {
	UserID   int64
	Sum      float64
	Category string
	Date     time.Time
}
