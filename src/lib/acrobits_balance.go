package acrobitsbalance

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	Version = "0.0.3"

	envPrefix   = "ACROBITS_BALANCE_"
	EnvPath     = envPrefix + "PATH"
	EnvAddr     = envPrefix + "ADDR"
	EnvCurrency = envPrefix + "CURRENCY"

	DefaultPath     = "/acrobits/balance"
	DefaultAddr     = "127.0.0.1:8080"
	DefaultCurrency = "USD"

	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 5 * time.Second
	DefaultIdleTimeout  = 30 * time.Second

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

type Config struct {
	Path         string                                `json:"path"`
	Addr         string                                `json:"addr"`
	Currency     string                                `json:"currency"`
	CertFile     string                                `json:"cert_file"`
	KeyFile      string                                `json:"key_file"`
	ReadTimeout  *time.Duration                        `json:"read_timeout"`
	WriteTimeout *time.Duration                        `json:"write_timeout"`
	IdleTimeout  *time.Duration                        `json:"idle_timeout"`
	Func         func(string, string) (float64, error) `json:"-"`
	path         string
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
		balance, err := c.Func(username, password)
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

func ListenAndServe(c Config) error {
	http.HandleFunc(c.Path, makeHandleFunc(&c))
	s := &http.Server{
		Addr:           c.Addr,
		ReadTimeout:    *c.ReadTimeout,
		WriteTimeout:   *c.WriteTimeout,
		IdleTimeout:    *c.IdleTimeout,
		MaxHeaderBytes: 1 << 12,
	}
	if c.CertFile != "" && c.KeyFile != "" {
		return s.ListenAndServeTLS(c.CertFile, c.KeyFile)
	}
	return s.ListenAndServe()
}
