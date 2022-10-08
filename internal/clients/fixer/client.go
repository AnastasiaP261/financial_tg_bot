package fixer

import (
	"encoding/json"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/exchange_rates"
	"io/ioutil"
	"net/http"
	"net/url"
)

// TokenGetter геттер токена апи
type TokenGetter interface {
	FixerAPIToken() string
}

// Client fixer клиент
type Client struct {
	tokenGetter TokenGetter
}

// New создать инстанс клиента fixer
func New(tokenGetter TokenGetter) *Client {
	return &Client{tokenGetter: tokenGetter}
}

// Response формат ответа
type Response struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
}

// GetExchangeRateToRUB получить валютные курсы для USD, EUR и CNY к RUB
func (c *Client) GetExchangeRateToRUB() (model.Response, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.apilayer.com",
		Path:   "fixer/latest",
		RawQuery: url.Values{
			"base":    []string{"RUB"},
			"symbols": []string{"USD", "EUR", "CNY"},
		}.Encode(),
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("apikey", c.tokenGetter.FixerAPIToken())
	if err != nil {
		return model.Response{}, errors.Wrap(err, "http.NewRequest")
	}

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return model.Response{}, errors.Wrap(err, "ioutil.ReadAll")
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return model.Response{}, errors.Wrap(err, "unmarshalling error")
	}

	return model.Response(response), nil
}
