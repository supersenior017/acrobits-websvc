package acrobitsbalance

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	Version = "0.1.2"

	envPrefix   = "ACROBITS_BALANCE_"
	EnvPath     = envPrefix + "PATH"
	EnvAddr     = envPrefix + "ADDR"
	EnvCurrency = envPrefix + "CURRENCY"

	DefaultPath     = "/acrobits/balance"
	DefaultAddr     = "127.0.0.1:8080"
	DefaultCurrency = "USD"

	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = DefaultReadTimeout
	DefaultIdleTimeout    = 30 * time.Second
	DefaultHandlerTimeout = DefaultIdleTimeout
)

type responseOK struct {
	BalanceString string  `json:"balanceString"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
}

func writeResponseOK(
	w http.ResponseWriter,
	balance float64,
	currency string,
) error {
	return json.NewEncoder(w).Encode(&responseOK{
		Balance:  balance,
		Currency: currency,
		BalanceString: currency + " " +
			strconv.FormatFloat(balance, 'G', -1, 64),
	})
}

type responseError struct {
	Message string `json:"message"`
}

func writeResponseError(w http.ResponseWriter, msg string) error {
	return json.NewEncoder(w).Encode(&responseError{msg})
}

type GetBalance func(context.Context, string, string) (float64, error)

type Config struct {
	Path           string         `json:"path"`
	Addr           string         `json:"addr"`
	Currency       string         `json:"currency"`
	CertFile       string         `json:"cert_file"`
	KeyFile        string         `json:"key_file"`
	ReadTimeout    *time.Duration `json:"read_timeout"`
	WriteTimeout   *time.Duration `json:"write_timeout"`
	IdleTimeout    *time.Duration `json:"idle_timeout"`
	HandlerTimeout *time.Duration `json:"handler_timeout"`
	GetBalance     GetBalance     `json:"-"`
	isSet          bool
}

func (c *Config) String() string {
	return ""
}

func (c *Config) Set(s string) error {
	var file *os.File
	var err error
	if s == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(s)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	if err = json.NewDecoder(file).Decode(c); err != nil {
		return err
	}
	c.SetDefaults()
	c.isSet = true
	return nil
}

func (c *Config) SetDefaults() {
	getenv := func(s, def string) string {
		if res := os.Getenv(s); res != "" {
			return res
		}
		return def
	}
	if c.Path == "" {
		c.Path = getenv(EnvPath, DefaultPath)
	}
	if c.Addr == "" {
		c.Addr = getenv(EnvAddr, DefaultAddr)
	}
	if c.Currency == "" {
		c.Currency = getenv(EnvCurrency, DefaultCurrency)
	}
	if c.ReadTimeout == nil {
		readTimeout := DefaultReadTimeout
		c.ReadTimeout = &readTimeout
	}
	if c.WriteTimeout == nil {
		writeTimeout := DefaultWriteTimeout
		c.WriteTimeout = &writeTimeout
	}
	if c.IdleTimeout == nil {
		idleTimeout := DefaultIdleTimeout
		c.IdleTimeout = &idleTimeout
	}
	if c.HandlerTimeout == nil {
		handlerTimeout := DefaultHandlerTimeout
		c.HandlerTimeout = &handlerTimeout
	}
}

func (c *Config) IsSet() bool {
	return c.isSet
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	writeResponseError(w, msg)
}

func makeHandleFunc(c *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		username := q.Get("username")
		password := q.Get("password")
		if username == "" || password == "" {
			code := http.StatusForbidden
			httpError(w, http.StatusText(code), code)
			return
		}
		balance, err := c.GetBalance(r.Context(), username, password)
		if err != nil {
			httpError(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		writeResponseOK(w, balance, c.Currency)
	}
}

func ListenAndServe(ctx context.Context, c Config) error {
	http.HandleFunc(c.Path, makeHandleFunc(&c))
	srv := &http.Server{
		Addr:           c.Addr,
		ReadTimeout:    *c.ReadTimeout,
		WriteTimeout:   *c.WriteTimeout,
		IdleTimeout:    *c.IdleTimeout,
		MaxHeaderBytes: 1 << 12,
		Handler: http.TimeoutHandler(
			http.DefaultServeMux,
			*c.HandlerTimeout,
			"",
		),
	}
	errchan := make(chan error, 1)
	go func() {
		if c.CertFile != "" && c.KeyFile != "" {
			errchan <- srv.ListenAndServeTLS(c.CertFile, c.KeyFile)
		} else {
			errchan <- srv.ListenAndServe()
		}
	}()
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(
			context.Background(),
			*c.HandlerTimeout,
		)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-errchan:
		return err
	}
}
