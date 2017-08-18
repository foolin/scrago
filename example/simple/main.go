package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"github.com/foolin/scrago"
	"log"
	"encoding/json"
	"os"
)

type ExampModel struct {
	Title string `scrago:"title"`
	Name string `scrago:"#main>.intro>h2::text()"`
	Description string `scrago:"#main>.intro>p::html()"`
	Intro string  `scrago:"#main>.intro::outerHtml()"`
	Keywords []string  `scrago:"#main .keywords::GetMyKeywords()"`
}

func (e *ExampModel) GetMyKeywords(s *goquery.Selection) ([]string, error) {
	v := s.Text()
	if v == ""{
		return nil, fmt.Errorf("not found keywords!")
	}
	arr := strings.Split(v, ",")
	for i := 0; i < len(arr); i++{
		arr[i] = strings.TrimSpace(arr[i])
	}
	return arr, nil
}

func main()  {
	examp := ExampModel{}
	s := scrago.New()
	err := s.HttpGetParser("https://raw.githubusercontent.com/foolin/scrago/master/example/data/example.html", &examp)
	if err != nil {
		log.Fatal(err)
	}else{
		printjson(examp)
	}
}

func printjson(v interface{})  {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	enc.Encode(v)
}
