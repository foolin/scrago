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
	Title string `css:"title"`
	Name string `css:"#main>.intro>h2::text()"`
	Description string `css:"#main>.intro>p::html()"`
	Keywords []string  `css:"#main .keywords::GetMyKeywords()"`
	Intro string  `css:"#main>.intro::outerHtml()"`
	TypeList []ExampTypeModel  `css:".typelist>ul>li"`
	SubModel ExampSubModel `css:".typelist"`
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
	TBool bool `css:"li[data-type='bool']"`
	TInt int `css:"li[data-type='int']"`
	TFloat float32 `css:"li[data-type='float']"`
	TString string `css:"li[data-type='string']"`
	TArray []string `css:"li[data-type='array'] li"`
}

type ExampTypeModel struct {
	Type string `css:"::attr(data-type)"`
	Value string  `css:"::text()"`
}

func main()  {
	examp := &ExampModel{}
	//examp := &ExampSubModel{}
	htmlContent, err := ioutil.ReadFile("./data/example.html")
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
