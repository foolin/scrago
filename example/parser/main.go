package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"github.com/foolin/scrago"
	"log"
	"encoding/json"
	"os"
	"io/ioutil"
)

type ExampModel struct {
	Title string `scrago:"title"`
	Name string `scrago:"#main>.intro>h2::text()"`
	Description string `scrago:"#main>.intro>p::html()"`
	Keywords []string  `scrago:"#main .keywords::GetMyKeywords()"`
	Intro string  `scrago:"#main>.intro::outerHtml()"`
	TypeList []ExampTypeModel  `scrago:".typelist>ul>li"`
	SubModel ExampSubModel `scrago:".typelist"`
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

type ExampSubModel struct {
	TBool bool `scrago:"li[data-type='bool']"`
	TInt int `scrago:"li[data-type='int']"`
	TFloat float32 `scrago:"li[data-type='float']"`
	TString string `scrago:"li[data-type='string']"`
	TArray []string `scrago:"li[data-type='array'] li"`
}

type ExampTypeModel struct {
	Type string `scrago:"::attr(data-type)"`
	Value string  `scrago:"::text()"`
}

func main()  {
	examp := &ExampModel{}
	//examp := &ExampSubModel{}
	htmlContent, err := ioutil.ReadFile("../data/example.html")
	if err != nil {
		log.Fatal(err)
	}
	document, _ := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
	err = scrago.ParserDocument(examp, document)
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
