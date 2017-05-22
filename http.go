package crawler

import (
	"net/http"
	"time"
	"net"
	"io/ioutil"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"io"
)

const UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"

func NewCrawler() *Crawler {
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
	return &Crawler{
		Client: client,
		UserAgent: UserAgent,
	}
}

func (c *Crawler) NewHttpRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	//request
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	return req, err
}

func (c *Crawler) HttpGetResponse(url string) (*http.Response, error) {
	return c.NewHttpRequest(http.MethodGet, url, nil)
}

func (c *Crawler) HttpPostResponse(url string, data string) (*http.Response, error) {
	//request
	req, err := c.NewHttpRequest(http.MethodPost, url, bytes.NewBufferString(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	//client
	return c.Do(req)
}

func (c *Crawler) HttpGetRaw(url string) ([]byte, error) {
	resp, err := c.HttpGetResponse(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//read body
	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return byteBody, nil
}

func (c *Crawler) HttpPostRaw(url string, data string) ([]byte, error) {
	//client
	resp, err := c.HttpPostResponse(url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//read body
	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return byteBody, nil
}

func (c *Crawler) HttpPostGizpRaw(url string, data []byte) ([]byte, error) {
	//zip compress
	var dataBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&dataBuffer)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}
	//request
	req, err := c.NewHttpRequest(http.MethodPost, url, &dataBuffer)
	req.Header.Add("Accept-Encoding", "gzip")
	if err != nil {
		return nil, err
	}
	//client
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//gzip umcompress
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	//read body
	byteBody, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return byteBody, nil
}

func (c *Crawler)  HttpSave(url string, file string) error {
	data, err := c.HttpGetRaw(url)
	if err != nil {
		return err
	}
	//mkdir
	err = os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}
	//write file
	err = ioutil.WriteFile(file, data, 0755)
	if err != nil {
		return err
	}
	return nil
}