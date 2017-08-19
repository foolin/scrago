# scrago

Scrago is an simpe, fast, extensible crawl page framework for golang.


# Install

```
 go get github.com/foolin/scrago
```

# Document

[Godoc](https://godoc.org/github.com/foolin/scrago "go document")

# Exmaple

Target page：
```html
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
</html>
```


### Step 1：
```go

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

```

### Step 2:
```go

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

```

### Step 3:
Execute result：

```json

{
    "Title": "Scrago exmaples",
    "Name": "Scrago framework",
    "Description": "An open source and collaborative framework for extracting the data you need from websites.\n            In a <b>fast</b>, <b>simple</b>, yet extensible way.",
    "Intro": "<div class=\"intro\">\n        <h2>Scrago framework</h2>\n        <p>An open source and collaborative framework for extracting the data you need from websites.\n            In a <b>fast</b>, <b>simple</b>, yet extensible way.</p>\n        <div class=\"keywords\">Scrago, Scrap, Spider, Crawl, GoLang, Simple, Easy</div>\n    </div>",
    "Keywords": [
        "Scrago",
        "Scrap",
        "Spider",
        "Crawl",
        "GoLang",
        "Simple",
        "Easy"
    ]
}

```

# Struct tag
Between selector and function use "::" symbol segmentation
```go
`scrago:"selector::function"`

```
* selector:
  Css selector, sea more：github.com/PuerkitoBio/goquery

* function:
  Get data function，default is text()。

  1.Inner function：
  - text() get text value.
  - html() get html vlaue.
  - outerHtml() get outer html value.
  - attr(xxx) get attribute value, eg：attr(href)。

  2.Write custom function：
```go

func (e *ExampModel) MyFunc(s *goquery.Selection) (MyReturnType, error) {
    //todo
    return ReturnValue, nil
}

```

   eg：
```go

type ExampModel struct {
    TextField string `scrago:"#xxx"`
    TextField2 string `scrago:".xxx::text()"`
    Link string `scrago:"a::attr(href)"`
    MyField string  `scrago:"#xxx::MyFunc()"`
}

func (e *ExampModel) MyFunc(s *goquery.Selection) (String, error) {
    //todo
    return s.Text(), nil
}

```


# Exmaples
 * [Simple](https://github.com/foolin/scrago/tree/master/example/simple "Simple Example")
 * [Parser](https://github.com/foolin/scrago/tree/master/example/parser "Parser Example")
 * [Quotesbot](https://github.com/foolin/scrago/tree/master/example/quotesbot "Quotesbot Example")

# Relative
 * github.com/PuerkitoBio/goquery