package websvc

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	envBalancePrefix   = envPrefix + "BALANCE_"
	EnvBalancePath     = envBalancePrefix + "PATH"
	EnvBalanceCurrency = envBalancePrefix + "CURRENCY"

	DefaultBalancePath     = "balance"
	DefaultBalanceCurrency = "USD"
)

type balanceResponse struct {
	BalanceString string  `json:"balanceString"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
}

func writeBalanceResponse(w http.ResponseWriter, balance Balance) error {
	return json.NewEncoder(w).Encode(&balanceResponse{
		Balance:  balance.Balance,
		Currency: balance.Currency,
		BalanceString: balance.Currency + " " +
			strconv.FormatFloat(balance.Balance, 'G', -1, 64),
	})
}

type Balance struct {
	Balance  float64
	Currency string
}

type BalanceFunc func(context.Context, Account) (Balance, error)

type BalanceConfig struct {
	Path     string      `json:"path"`
	Currency string      `json:"currency"`
	Func     BalanceFunc `json:"-"`
}

func (c *BalanceConfig) SetDefaults() {
	if c.Path == "" {
		c.Path = getenv(EnvBalancePath, DefaultBalancePath)
	}
	if c.Currency == "" {
		c.Currency = getenv(EnvBalanceCurrency, DefaultBalanceCurrency)
	}
}

func makeBalanceHandleFunc(c *BalanceConfig) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		username := q.Get("username")
		password := q.Get("password")
		if username == "" && password == "" {
			code := http.StatusBadRequest
			httpError(w, http.StatusText(code), code)
			return
		}
		balance, err := c.Func(r.Context(), Account{username, password})
		if err != nil {
			httpError(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		if balance.Currency == "" {
			balance.Currency = c.Currency
		}
		writeBalanceResponse(w, balance)
	}
}
