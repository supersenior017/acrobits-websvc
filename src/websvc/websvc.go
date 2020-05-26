package websvc

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const (
	Version = "0.1.4"

	envPrefix = "ACROBITS_WEBSVC_"
	EnvPath   = envPrefix + "PATH"
	EnvAddr   = envPrefix + "ADDR"

	DefaultPath = "/acrobits/"
	DefaultAddr = "127.0.0.1:8080"

	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = DefaultReadTimeout
	DefaultIdleTimeout    = 30 * time.Second
	DefaultHandlerTimeout = DefaultIdleTimeout
)

type Account struct {
	Username string
	Password string
}

type responseError struct {
	Message string `json:"message"`
}

func writeResponseError(w http.ResponseWriter, msg string) error {
	return json.NewEncoder(w).Encode(&responseError{msg})
}

func setConstDurationRef(t **time.Duration, c time.Duration) {
	*t = &c
}

type Config struct {
	Path           string         `json:"path"`
	Addr           string         `json:"addr"`
	Balance        BalanceConfig  `json:"balance"`
	CertFile       string         `json:"cert_file"`
	KeyFile        string         `json:"key_file"`
	ReadTimeout    *time.Duration `json:"read_timeout"`
	WriteTimeout   *time.Duration `json:"write_timeout"`
	IdleTimeout    *time.Duration `json:"idle_timeout"`
	HandlerTimeout *time.Duration `json:"handler_timeout"`
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
	if c.Path == "" {
		c.Path = getenv(EnvPath, DefaultPath)
	}
	if c.Addr == "" {
		c.Addr = getenv(EnvAddr, DefaultAddr)
	}
	c.Balance.SetDefaults()
	if c.ReadTimeout == nil {
		setConstDurationRef(&c.ReadTimeout, DefaultReadTimeout)
	}
	if c.WriteTimeout == nil {
		setConstDurationRef(&c.WriteTimeout, DefaultWriteTimeout)
	}
	if c.IdleTimeout == nil {
		setConstDurationRef(&c.IdleTimeout, DefaultIdleTimeout)
	}
	if c.HandlerTimeout == nil {
		setConstDurationRef(&c.HandlerTimeout, DefaultHandlerTimeout)
	}
}

func (c *Config) IsSet() bool {
	return c.isSet
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	writeResponseError(w, msg)
}

func ListenAndServe(ctx context.Context, c Config) error {
	http.HandleFunc(
		urljoin(c.Path, c.Balance.Path),
		makeBalanceHandleFunc(&c.Balance),
	)
	timeoutMsg, _ := json.Marshal(&responseError{"timeout"})
	srv := &http.Server{
		Addr:           c.Addr,
		ReadTimeout:    *c.ReadTimeout,
		WriteTimeout:   *c.WriteTimeout,
		IdleTimeout:    *c.IdleTimeout,
		MaxHeaderBytes: 1 << 12,
		Handler: &jsonHandler{http.TimeoutHandler(
			http.DefaultServeMux,
			*c.HandlerTimeout,
			string(timeoutMsg),
		)},
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
