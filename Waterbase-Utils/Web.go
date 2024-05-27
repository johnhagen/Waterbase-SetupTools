package WBU

import (
	"net/http"
	"time"
)

func WebClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     1 * time.Second,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 1,
	}

	return &http.Client{Transport: transport}
}
