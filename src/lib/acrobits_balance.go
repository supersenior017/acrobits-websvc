package acrobitsbalance

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	Version = "0.0.5"

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

	ArgStdin = "-"
)

type xmlResponse struct {
	XMLName       xml.Name `xml:"response"`
	BalanceString string   `xml:"balanceString"`
	Balance       float64  `xml:"balance"`
	Currency      string   `xml:"currency"`
}

func newXmlResponse(balance float64, currency string) *xmlResponse {
	r := &xmlResponse{}
	r.BalanceString = currency + " " +
		strconv.FormatFloat(balance, 'G', -1, 64)
	r.Balance = balance
	r.Currency = currency
	return r
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
	Func           GetBalance     `json:"-"`
	path           string
}

func (c Config) String() string {
	return ""
}

func (c *Config) Set(s string) error {
	var file *os.File
	var err error
	if s == ArgStdin {
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
	c.path = s
	return nil
}

func (c *Config) SetDefaults() {
	getenv := func(s, def string) string {
		res := os.Getenv(s)
		if res == "" {
			return def
		}
		return res
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

func (c Config) FilePath() string {
	return c.path
}

func makeHandleFunc(c *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		username := q.Get("username")
		password := q.Get("password")
		if username == "" || password == "" {
			http.Error(
				w,
				http.StatusText(http.StatusForbidden),
				http.StatusForbidden,
			)
			return
		}
		balance, err := c.Func(r.Context(), username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(
			newXmlResponse(balance, c.Currency),
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ListenAndServe(ctx context.Context, c Config) error {
	http.HandleFunc(c.Path, makeHandleFunc(&c))
	s := &http.Server{
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
			errchan <- s.ListenAndServeTLS(c.CertFile, c.KeyFile)
		} else {
			errchan <- s.ListenAndServe()
		}
	}()
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(
			context.Background(),
			*c.HandlerTimeout,
		)
		defer cancel()
		return s.Shutdown(ctx)
	case err := <-errchan:
		return err
	}
}
