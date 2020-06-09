package websvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	Version = "0.1.10"

	envPrefix   = "ACROBITS_WEBSVC_"
	EnvPath     = envPrefix + "PATH"
	EnvAddr     = envPrefix + "ADDR"
	EnvCertFile = envPrefix + "CERT_FILE"
	EnvKeyFile  = envPrefix + "KEY_FILE"

	EnvReadTimeout    = envPrefix + "READ_TIMEOUT"
	EnvWriteTimeout   = envPrefix + "WRITE_TIMEOUT"
	EnvIdleTimeout    = envPrefix + "IDLE_TIMEOUT"
	EnvHandlerTimeout = envPrefix + "HANDLER_TIMEOUT"

	DefaultPath = "/acrobits/"
	DefaultAddr = "127.0.0.1:8080"

	DefaultReadTimeout    = 1 << 2 * time.Second
	DefaultWriteTimeout   = DefaultReadTimeout
	DefaultIdleTimeout    = 1 << 5 * time.Second
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

type Duration time.Duration

func (d *Duration) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.Set(s)
}

func (d Duration) Unwrap() time.Duration {
	return time.Duration(d)
}

func setConfigEnabled(dest *bool, envKey string) error {
	if val := os.Getenv(envKey); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		*dest = b
	}
	return nil
}

func setConfigTimeout(
	dest **Duration,
	envKey string,
	defaultValue time.Duration,
) error {
	if val := os.Getenv(envKey); val != "" {
		var d Duration
		if err := d.Set(val); err != nil {
			return err
		}
		*dest = &d
	} else if *dest == nil {
		d := Duration(defaultValue)
		*dest = &d
	}
	return nil
}

func setConfigString(dest *string, envKey, defaultValue string) {
	if val := os.Getenv(envKey); val != "" {
		*dest = val
	} else if *dest == "" {
		*dest = defaultValue
	}
}

type Config struct {
	Path           string        `json:"path"`
	Addr           string        `json:"addr"`
	Balance        BalanceConfig `json:"balance"`
	Rate           RateConfig    `json:"rate"`
	CertFile       string        `json:"cert_file"`
	KeyFile        string        `json:"key_file"`
	ReadTimeout    *Duration     `json:"read_timeout"`
	WriteTimeout   *Duration     `json:"write_timeout"`
	IdleTimeout    *Duration     `json:"idle_timeout"`
	HandlerTimeout *Duration     `json:"handler_timeout"`
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
	return c.SetDefaults()
}

func (c *Config) SetDefaults() error {
	setConfigString(&c.Path, EnvPath, DefaultPath)
	setConfigString(&c.Addr, EnvAddr, DefaultAddr)
	if err := c.Balance.SetDefaults(); err != nil {
		return err
	}
	if err := c.Rate.SetDefaults(); err != nil {
		return err
	}
	if c.CertFile == "" {
		c.CertFile = os.Getenv(EnvCertFile)
	}
	if c.KeyFile == "" {
		c.KeyFile = os.Getenv(EnvKeyFile)
	}
	if err := setConfigTimeout(
		&c.ReadTimeout,
		EnvReadTimeout,
		DefaultReadTimeout,
	); err != nil {
		return err
	}
	if err := setConfigTimeout(
		&c.WriteTimeout,
		EnvWriteTimeout,
		DefaultWriteTimeout,
	); err != nil {
		return err
	}
	if err := setConfigTimeout(
		&c.IdleTimeout,
		EnvIdleTimeout,
		DefaultIdleTimeout,
	); err != nil {
		return err
	}
	if err := setConfigTimeout(
		&c.HandlerTimeout,
		EnvHandlerTimeout,
		DefaultHandlerTimeout,
	); err != nil {
		return err
	}
	c.isSet = true
	return nil
}

func (c *Config) IsSet() bool {
	return c.isSet
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	writeResponseError(w, msg)
}

func addHandlers(m *http.ServeMux, c *Config) error {
	hasEnabled := false
	if c.Balance.Enabled {
		if c.Balance.Func == nil {
			return errors.New("balance func is nil")
		}
		m.HandleFunc(
			urlMustJoin(c.Path, c.Balance.Path),
			makeBalanceHandleFunc(c.Balance),
		)
		hasEnabled = true
	}
	if c.Rate.Enabled {
		if c.Rate.Func == nil {
			return errors.New("rate func is nil")
		}
		m.HandleFunc(
			urlMustJoin(c.Path, c.Rate.Path),
			makeRateHandleFunc(c.Rate),
		)
		hasEnabled = true
	}
	if !hasEnabled {
		return errors.New("no enabled modules")
	}
	return nil
}

func ListenAndServe(ctx context.Context, c Config) error {
	if !c.isSet {
		if err := c.SetDefaults(); err != nil {
			return err
		}
	}
	mux := http.NewServeMux()
	if err := addHandlers(mux, &c); err != nil {
		return err
	}
	timeoutMsg, _ := json.Marshal(&responseError{"timeout"})
	handlerTimeout := c.HandlerTimeout.Unwrap()
	srv := &http.Server{
		Addr:           c.Addr,
		ReadTimeout:    c.ReadTimeout.Unwrap(),
		WriteTimeout:   c.WriteTimeout.Unwrap(),
		IdleTimeout:    c.IdleTimeout.Unwrap(),
		MaxHeaderBytes: 1 << 12,
		Handler: &jsonResponseHandler{http.TimeoutHandler(
			mux,
			handlerTimeout,
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
			handlerTimeout,
		)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-errchan:
		return err
	}
}
