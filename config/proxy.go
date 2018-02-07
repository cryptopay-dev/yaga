package config

import (
	"net/http"
	"net/url"
)

// Proxy default configuration
type Proxy string

// String return proxy-url
func (p Proxy) String() string {
	return string(p)
}

// Apply proxy
func (p Proxy) Apply(t http.RoundTripper) error {
	var (
		err      error
		proxyURL *url.URL
		strURL   = p.String()
	)

	// Not need if empty URL:
	if len(strURL) == 0 {
		return nil
	}

	// Try to parse url..
	if proxyURL, err = url.Parse(p.String()); err != nil {
		return err
	}

	// Set default transport:
	t = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	return nil
}
