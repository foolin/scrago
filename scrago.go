package scrago

import (
	"net/http"
	"net"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"

type Scrago struct {
	*http.Client
	UserAgent string
}

func New() *Scrago {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 15 * time.Second,
	}
	return &Scrago{
		Client: client,
		UserAgent: UserAgent,
	}
}