package acrobitswebsvc

import (
	"net/http"
	"net/url"
	"os"
)

func getenv(s, def string) string {
	if v := os.Getenv(s); v != "" {
		return v
	}
	return def
}

func urljoin(base, ref string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	refURL, err := url.Parse(ref)
	if err != nil {
		panic(err)
	}
	return baseURL.ResolveReference(refURL).String()
}

type jsonHandler struct {
	handler http.Handler
}

func (h *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.handler.ServeHTTP(w, r)
}
