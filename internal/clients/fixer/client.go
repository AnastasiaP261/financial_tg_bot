package fixer

import (
	"context"
	"encoding/json"
	"io"
	"log"
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
		log.Println("[FIXER]: Запущен процесс получения актуальных курсов валют из fixer")
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
		log.Printf("[FIXER CLIENT ERR]: %s\n", errors.Wrap(err, "http.NewRequest").Error())
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("[FIXER CLIENT ERR]: %s\n", errors.Wrap(err, "client.Do").Error())
		return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[FIXER CLIENT ERR]: %s\n", errors.Wrap(err, "ioutil.ReadAll").Error())
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("[FIXER CLIENT ERR]: %s\n", errors.Wrap(err, "unmarshalling error").Error())
	}

	if response.Success {
		log.Printf("[FIXER]: получен текущий курс: %+v\n", response.Rates)
		c.dataAccessorWrite(response.Rates)
	} else {
		log.Println("[FIXER CLIENT ERR]: client request failed")
	}
}

// GetExchangeRateToRUB получить валютные курсы для USD, EUR и CNY к RUB
func (c *Client) GetExchangeRateToRUB() map[string]float64 {
	return c.dataAccessorRead()
}
