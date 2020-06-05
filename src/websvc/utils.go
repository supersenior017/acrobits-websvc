package websvc

import (
	"net/http"
	"net/url"
	"strconv"
)

func urljoin(base, ref string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	refURL, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(refURL).String(), nil
}

func urlMustJoin(base, ref string) string {
	res, err := urljoin(base, ref)
	if err != nil {
		panic(err)
	}
	return res
}

type jsonHandler struct {
	handler http.Handler
}

func (h *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.handler.ServeHTTP(w, r)
}

func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
