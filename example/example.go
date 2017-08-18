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
	document, _ := goquery.NewDocumentFromReader(strings.NewReader(exampContent))
	err := scrago.ParserDocument(examp, document)
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

var exampContent = `
<!doctype html>
<html class="no-js" lang="">

<head>
    <meta charset="utf-8">
    <title>Scrago exmaples</title>
</head>

<body>
<div id="header">
    <div class="container">
        <div class="clearfix">
            <div class="logo">
                <a href="https://github.com/foolin/scrago" title="Scrago exmaple">
                    <h1 title="Scrago exmaple - crawl framework for go">Scrago exmaple</h1>
                </a>
            </div>
        </div>
    </div>
</div>

<div class="navlink">
    <div class="container">
        <ul class="clearfix">
            <li ><a href="/">Index</a></li>
            <li ><a href="/list/web" title="web site">Web page</a></li>
            <li ><a href="/list/pc" title="pc page">Pc Page</a></li>
            <li ><a href="/list/mobile" title="mobile page">Mobile Page</a></li>
        </ul>
    </div>
</div>

<div id="main">
	<div class="intro">
		<h2>Scrago framework</h2>
		<p>An open source and collaborative framework for extracting the data you need from websites.
	In a <b>fast</b>, <b>simple</b>, yet extensible way.</p>
		<div class="keywords">Scrago, Scrap, Spider, Crawl, GoLang, Simple, Easy</div>
	</div>
	<div class="typelist">
		<ul>
			<li data-type="bool">true</li>
			<li data-type="int">123</li>
			<li data-type="float">45.6</li>
			<li data-type="string">hello</li>
			<li data-type="array">
				<ol>
					<li>Aa</li>
					<li>Bb</li>
					<li>Cc</li>
				</ol>
			</li>
		</ul>
	</div>

</div>

</body>
</html>`
