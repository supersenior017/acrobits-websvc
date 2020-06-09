package websvc

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	envBalancePrefix   = envPrefix + "BALANCE_"
	EnvBalancePath     = envBalancePrefix + "PATH"
	EnvBalanceCurrency = envBalancePrefix + "CURRENCY"
	EnvBalanceEnabled  = envBalancePrefix + "ENABLED"

	DefaultBalancePath     = "balance"
	DefaultBalanceCurrency = "USD"
)

type balanceResponse struct {
	BalanceString string  `json:"balanceString"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
}

func writeBalanceResponse(w http.ResponseWriter, b Balance) error {
	return json.NewEncoder(w).Encode(&balanceResponse{
		Balance:       b.Balance,
		Currency:      b.Currency,
		BalanceString: b.Currency + " " + ftoa(b.Balance),
	})
}

type Balance struct {
	Balance  float64
	Currency string
}

func (b *Balance) SetDefaults(c *BalanceConfig) {
	if b.Currency == "" {
		b.Currency = c.Currency
	}
}

type BalanceFunc func(context.Context, Account) (Balance, error)

type BalanceConfig struct {
	Path     string      `json:"path"`
	Currency string      `json:"currency"`
	Enabled  bool        `json:"enabled"`
	Func     BalanceFunc `json:"-"`
}

func (c *BalanceConfig) SetDefaults() error {
	setConfigString(&c.Path, EnvBalancePath, DefaultBalancePath)
	setConfigString(&c.Currency, EnvBalanceCurrency, DefaultBalanceCurrency)
	return setConfigEnabled(&c.Enabled, EnvBalanceEnabled)
}

func makeBalanceHandleFunc(c BalanceConfig) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		account := Account{
			Username: q.Get("username"),
			Password: q.Get("password"),
		}
		if account.Username == "" && account.Password == "" {
			code := http.StatusBadRequest
			httpError(w, http.StatusText(code), code)
			return
		}
		balance, err := c.Func(r.Context(), account)
		if err != nil {
			httpError(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		balance.SetDefaults(&c)
		writeBalanceResponse(w, balance)
	}
}
