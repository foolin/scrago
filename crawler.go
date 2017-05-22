package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type Crawler struct {
	*http.Client
	UserAgent string
}

func (c *Crawler) HttpGetDoc(url string) (*goquery.Document, error) {
	resp, err := c.HttpGetResponse(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}

func (c *Crawler) HttpPostDoc(url string, data string) (*goquery.Document, error) {
	resp, err := c.HttpPostResponse(url, data)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}
