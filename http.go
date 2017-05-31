package scrago

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"io"
)


func (c *Scrago) NewHttpRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	//request
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	return req, err
}

func (c *Scrago) HttpGetResponse(url string) (*http.Response, error) {
	req, err := c.NewHttpRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//client
	return c.Do(req)
}

func (c *Scrago) HttpPostResponse(url string, data string) (*http.Response, error) {
	//request
	req, err := c.NewHttpRequest(http.MethodPost, url, bytes.NewBufferString(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	//client
	return c.Do(req)
}

func (c *Scrago) HttpGetRaw(url string) ([]byte, error) {
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

func (c *Scrago) HttpPostRaw(url string, data string) ([]byte, error) {
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

func (c *Scrago) HttpPostGizpRaw(url string, data []byte) ([]byte, error) {
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

func (c *Scrago)  HttpSave(url string, file string) error {
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