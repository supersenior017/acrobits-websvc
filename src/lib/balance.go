package acrobitswebsvc

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

type responseBalance struct {
	BalanceString string  `json:"balanceString"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
}

func writeResponseBalance(
	w http.ResponseWriter,
	balance float64,
	currency string,
) error {
	return json.NewEncoder(w).Encode(&responseBalance{
		Balance:  balance,
		Currency: currency,
		BalanceString: currency + " " +
			strconv.FormatFloat(balance, 'G', -1, 64),
	})
}

type GetBalance func(context.Context, Params) (float64, error)

type Balance struct {
	Path     string     `json:"path"`
	Currency string     `json:"currency"`
	Func     GetBalance `json:"-"`
}

func (b *Balance) SetDefaults() {
	if b.Path == "" {
		b.Path = getenv(EnvBalancePath, DefaultBalancePath)
	}
	if b.Currency == "" {
		b.Currency = getenv(EnvBalanceCurrency, DefaultBalanceCurrency)
	}
}

func makeBalanceHandleFunc(b *Balance) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Context()
		q := r.URL.Query()
		username := q.Get("username")
		password := q.Get("password")
		if username == "" && password == "" {
			code := http.StatusBadRequest
			httpError(w, http.StatusText(code), code)
			return
		}
		balance, err := b.Func(r.Context(), Params{username, password})
		if err != nil {
			httpError(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeResponseBalance(w, balance, b.Currency)
	}
}
