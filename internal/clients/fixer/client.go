package fixer

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	sleepTime    = 60 * time.Minute
	fixerTimeout = 3 * time.Minute
)

// TokenGetter геттер токена апи
type TokenGetter interface {
	FixerAPIToken() string
}

// Client fixer клиент
type Client struct {
	tokenGetter TokenGetter

	exchangeRates map[string]float64
	mutex         sync.RWMutex
}

// New создать инстанс клиента fixer
func New(ctx context.Context, tokenGetter TokenGetter) *Client {
	cl := &Client{tokenGetter: tokenGetter}

	go cl.process(ctx)

	return cl
}

func (c *Client) process(ctx context.Context) {
	ticker := time.NewTicker(sleepTime)

	run := func() {
		logs.Info("fixer background process running")
		ctxNew, cancel := context.WithTimeout(ctx, fixerTimeout)
		c.getData(ctxNew)
		cancel()
	}

	run()

	select {
	case <-ctx.Done():
		ticker.Stop()
		return
	case <-ticker.C:
		run()
	}
}

func (c *Client) dataAccessorWrite(newData map[string]float64) {
	c.mutex.Lock()
	c.exchangeRates = newData
	c.mutex.Unlock()
}

func (c *Client) dataAccessorRead() map[string]float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.exchangeRates
}

// Response формат ответа fixer
type Response struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
}

// getData получает из апи данные по курсу валют на текущий момент времени
func (c *Client) getData(ctx context.Context) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.apilayer.com",
		Path:   "fixer/latest",
		RawQuery: url.Values{
			"access_key": []string{c.tokenGetter.FixerAPIToken()},
			"base":       []string{"RUB"},
			"symbols":    []string{"USD,EUR,CNY"},
		}.Encode(),
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req = req.WithContext(ctx)
	req.Header.Set("apikey", c.tokenGetter.FixerAPIToken())
	if err != nil {
		logs.Error(errors.Wrap(err, "http.NewRequest").Error())
		return
	}

	res, err := client.Do(req)
	if err != nil {
		logs.Error(errors.Wrap(err, "client.Do").Error())
		return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logs.Error(errors.Wrap(err, "ioutil.ReadAll").Error())
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		logs.Error(errors.Wrap(err, "unmarshalling error").Error())
	}

	if response.Success {
		logs.Info("current exchange rate received", zap.Any("rate", response.Rates))
		c.dataAccessorWrite(response.Rates)
	} else {
		logs.Error("client request failed")
	}
}

// getDataFromDate получает из апи данные по курсу валют на указанную дату
func (c *Client) getDataFromDate(ctx context.Context, y, m, d int) (map[string]float64, error) {
	if y <= 1999 {
		return nil, errors.New("invalid date")
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.apilayer.com",
		Path:   fmt.Sprintf("fixer/%d-%02d-%02d", y, m, d),
		RawQuery: url.Values{
			"access_key": []string{c.tokenGetter.FixerAPIToken()},
			"base":       []string{"RUB"},
			"symbols":    []string{"USD,EUR,CNY"},
		}.Encode(),
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req = req.WithContext(ctx)
	req.Header.Set("apikey", c.tokenGetter.FixerAPIToken())
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling error")
	}

	if response.Success {
		logs.Info(
			"current exchange rate for date received",
			zap.String("date", fmt.Sprintf("%02d.%02d.%d", d, m, y)),
			zap.Any("rate", response.Rates),
		)
		return response.Rates, nil
	} else {
		logs.Error("client request failed")
		return nil, errors.New("client request failed with unknown error")
	}
}

// GetExchangeRateToRUB получить последние валютные курсы для USD, EUR и CNY к RUB
func (c *Client) GetExchangeRateToRUB() map[string]float64 {
	return c.dataAccessorRead()
}

// GetExchangeRateToRUBFromDate получить валютные курсы для USD, EUR и CNY к RUB за конкретную дату
func (c *Client) GetExchangeRateToRUBFromDate(ctx context.Context, y, m, d int) (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(ctx, fixerTimeout)
	defer cancel()

	res, err := c.getDataFromDate(ctx, y, m, d)
	if err != nil {
		return nil, errors.Wrap(err, "getDataFromDate")
	}
	return res, nil
}
