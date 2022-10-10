package purchases

import (
	"github.com/pkg/errors"
	"strings"
)

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

func (m *Model) StrToCurrency(str string) (Currency, error) {
	str = strings.ToUpper(strings.TrimSpace(str))
	switch str {
	case "RUB":
		return RUB, nil
	case "USD":
		return USD, nil
	case "EUR":
		return EUR, nil
	case "CNY":
		return CNY, nil
	default:
		return 0, errors.New("invalid currency")
	}
}

func (m *Model) currencyToStr(cy Currency) (string, error) {
	switch cy {
	case RUB:
		return "RUB", nil
	case USD:
		return "USD", nil
	case EUR:
		return "EUR", nil
	case CNY:
		return "CNY", nil
	default:
		return "", errors.New("invalid currency")
	}
}

// RateToRUB курс валют к RUB
type RateToRUB struct {
	USD float64
	EUR float64
	CNY float64
}

// конвертирует сумму в валюте в рубли
func (m *Model) toRUB(userCurrency Currency, sum float64, rates RateToRUB) (float64, error) {
	switch userCurrency {
	case USD:
		return sum / rates.USD, nil
	case EUR:
		return sum / rates.EUR, nil
	case CNY:
		return sum / rates.CNY, nil
	case RUB:
		return sum, nil
	default:
		return 0, errors.New("invalid currency")
	}
}

// конвертирует сумму в рублях в указанную валюту
func (m *Model) rubToCurrentCurrency(userCurrency Currency, sum float64, rates RateToRUB) (float64, error) {
	switch userCurrency {
	case USD:
		return sum * rates.USD, nil
	case EUR:
		return sum * rates.EUR, nil
	case CNY:
		return sum * rates.CNY, nil
	case RUB:
		return sum, nil
	default:
		return 0, errors.New("invalid currency")
	}
}
