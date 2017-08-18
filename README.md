# scrago

Scrago is an open source and collaborative framework for extracting the data you need from websites.
In a fast, simple, yet extensible way.

# 中文文档

scrago是一个基于golang的爬虫框架，通过一种快速、简单、可扩展的方式，从网站中提取你需要的数据。


# 特点：
 * 快速
 * 简单
 * 可扩展

# 安装

```
 go get github.com/foolin/scrago
```

# 文档

[Document](https://godoc.org/github.com/foolin/scrago "go document")

# 示例

抓取目标页面：
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


### 第一步：
创建抓取数据struct，代码：
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

### 第二步:
编写抓取逻辑：
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

### 第三步:
执行并返回结果：

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

# Struct标签说明
tag使用scrago作为标签标示，选择器和方法之间用::分开，语法如下：
```go
`scrago:"selector::function"`

```
* selector:
  CSS选择器，类似jquery语法，具体使用请参考：github.com/PuerkitoBio/goquery

* function:
  函数方法，可自定义。如果省略，则默认是text方法。

  1.自带方法：
  - text 获取文本
  - html 获取html
  - outerHtml 获取整个节点html
  - attr(xxx) 获取节点属性，例如：attr(href)则获取<a href="http://www.liufu.me">liufu</a>中的href属性值。

  2.自定义方法：
  struct对象如下：
```go

func (e *ExampModel) 函数名(s *goquery.Selection) (返回类型, error) {
    //todo
    return 返回值, nil
}

```

  例如：
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


# 更多示例
 * [Simple](https://github.com/foolin/scrago/tree/master/example/simple "Simple Example")
 * [Parser](https://github.com/foolin/scrago/tree/master/example/parser "Parser Example")

# 依赖
 * github.com/PuerkitoBio/goquery