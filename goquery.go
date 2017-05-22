package scrago

import "github.com/PuerkitoBio/goquery"

func (c *Scrago) HttpGetDoc(url string) (*goquery.Document, error) {
	resp, err := c.HttpGetResponse(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}

func (c *Scrago) HttpPostDoc(url string, data string) (*goquery.Document, error) {
	resp, err := c.HttpPostResponse(url, data)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}
