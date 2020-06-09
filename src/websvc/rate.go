package websvc

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	envRatePrefix        = envPrefix + "RATE_"
	EnvRatePath          = envRatePrefix + "PATH"
	EnvRateCurrency      = envRatePrefix + "CURRENCY"
	EnvRateSpecification = envRatePrefix + "SPECIFICATION"
	EnvRateEnabled       = envRatePrefix + "ENABLED"

	DefaultRatePath          = "rate"
	DefaultRateCurrency      = "¢"
	DefaultRateSpecification = "min."
)

type rateResponse struct {
	CallRateString         string `json:"callRateString"`
	MessageRateString      string `json:"messageRateString"`
	SmartCallRateString    string `json:"smartCallRateString"`
	SmartMessageRateString string `json:"smartMessageRateString"`
}

func writeRateResponse(w http.ResponseWriter, r Rate) error {
	return json.NewEncoder(w).Encode(&rateResponse{
		CallRateString:         r.Call.format(r.Currency),
		MessageRateString:      ftoa(r.Message) + r.Currency,
		SmartCallRateString:    r.SmartCall.format(r.Currency),
		SmartMessageRateString: ftoa(r.SmartMessage) + r.Currency,
	})
}

type CallItem struct {
	Price         float64
	Specification string
}

func (c CallItem) format(currency string) string {
	return ftoa(c.Price) + currency + " " + c.Specification
}

type Rate struct {
	Call         CallItem
	Message      float64
	SmartCall    CallItem
	SmartMessage float64
	Currency     string
}

func (r *Rate) SetDefaults(c *RateConfig) {
	if r.Call.Specification == "" {
		r.Call.Specification = c.Specification
	}
	if r.SmartCall.Specification == "" {
		r.SmartCall.Specification = c.Specification
	}
	if r.Currency == "" {
		r.Currency = c.Currency
	}
}

type RateParams struct {
	Account      Account
	TargetNumber string
	SmartURI     string
}

type RateFunc func(context.Context, RateParams) (Rate, error)

type RateConfig struct {
	Path          string   `json:"path"`
	Currency      string   `json:"currency"`
	Specification string   `json:"specification"`
	Enabled       bool     `json:"enabled"`
	Func          RateFunc `json:"-"`
}

func (c *RateConfig) SetDefaults() error {
	setConfigString(&c.Path, EnvRatePath, DefaultRatePath)
	setConfigString(&c.Currency, EnvRateCurrency, DefaultRateCurrency)
	setConfigString(
		&c.Specification,
		EnvRateSpecification,
		DefaultRateSpecification,
	)
	return setConfigEnabled(&c.Enabled, EnvRateEnabled)
}

func makeRateHandleFunc(c RateConfig) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		params := RateParams{
			Account: Account{
				Username: q.Get("username"),
				Password: q.Get("password"),
			},
			TargetNumber: q.Get("targetNumber"),
			SmartURI:     q.Get("smartUri"),
		}
		if params.TargetNumber == "" && params.SmartURI == "" {
			code := http.StatusBadRequest
			httpError(w, http.StatusText(code), code)
			return
		}
		rate, err := c.Func(r.Context(), params)
		if err != nil {
			httpError(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		rate.SetDefaults(&c)
		writeRateResponse(w, rate)
	}
}
