package main

import (
	"github.com/foolin/scrago"
	"log"
	"os"
	"encoding/json"
)

type Quotesbot struct{
	List    []QuotesItem `scrago:"div.quote"`
	NextUrl string `scrago:"li.next > a::attr(href)"`
}

type QuotesItem struct {
	Text string `scrago:"span.text::text"`
	Author string `scrago:"small.author::text"`
	Tags []string `scrago:"div.tags > a.tag::text"`
}

func main()  {
	quot := Quotesbot{}
	s := scrago.New()
	err := s.HttpGetParser("http://quotes.toscrape.com/", &quot)
	if err != nil {
		log.Fatal(err)
	}else{
		printjson(quot)
	}
}

func printjson(v interface{})  {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	enc.Encode(v)
}
