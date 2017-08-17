package main

import (
	"reflect"
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"github.com/foolin/scrago"
)

//::text()
//::value()
//::html()
//::attr(xxx)
var rxFunc = regexp.MustCompile("^\\s*([a-zA-Z]+)\\s*\\(([^\\)]*)\\)\\s*$")



type DocModel struct{
	Name string `css:".header h1::attr(title)"`
	Head *SubModel `css:".header"`
}

func (d *DocModel) GetterTitle(s *goquery.Selection) (string, error) {
	return "ok", nil
}

type SubModel struct {
	Title string `css:"h1::attr(title)"`
}

func main()  {
	m := &DocModel{}
	//tv := reflect.ValueOf(m).MethodByName("GetterTitle")
	//log.Printf("%#v", reflect.ValueOf(m))
	//log.Printf("%#v", reflect.ValueOf(m).Elem())
	//log.Fatalf("%v", tv.IsValid())
	document, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	err := scrago.ParserField(m, document.Selection)
	if err != nil {
		log.Fatal(err)
	}else{
		log.Printf("value: %#v\n---\n%#v", m, m.Head)
	}
}

func main1()  {
	model := &DocModel{}
	runModel(model)

	log.Printf("%v", 1)

	selectors := strings.Split(`#id div[name='attr']::attr(aaa, aaaa, ccc)`, "::")
	log.Printf(selectors[0])
	log.Printf(selectors[1])
	log.Printf("%q\n====ddddds====\n", rxFunc.FindStringSubmatch(selectors[1]))

	document, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	fv, _ := goquery.OuterHtml(document.Find(".crumbs a"))
	log.Printf("%#v", fv)
	log.Printf("%#v", model)
	println("hello world")
}

func runModel(m interface{})  {
	rt := reflect.TypeOf(m)
	rv := reflect.ValueOf(m)
	if rt.Kind() == reflect.Ptr{
		rt = rt.Elem()
		rv = rv.Elem()
	}
	css := rt.Field(0).Tag.Get("css")
	rv.Field(0).SetString(css)
	log.Printf("css: %v", css)
	log.Printf("field type: %v", rt.Field(0).Type.Kind() == reflect.String)
}


var htmlContent = `
<!doctype html>
<html class="no-js" lang="">

<head>
    <meta charset="utf-8">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <title>
关于我们_七夕阅读网
</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="shortcut icon" href="/static/img/favicon.ico" type="image/x-icon" />
    <link rel="apple-touch-icon" href="/static/img/touch_icon.png">


    <link rel="stylesheet" href="/static/css/normalize.css">
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="stylesheet" href="/static/css/app.css">


</head>

<body>


<div class="header">
    <div class="container">
        <div class="clearfix">
            <div class="logo">
                <a href="/" title="七夕阅读网">
                    <h1 class="abc" data="aaaa" title="七夕阅读网-qixi.us">七夕阅读网</h1>
                </a>
            </div>
            <div class="banner">

            </div>
        </div>
    </div>
</div>

<div class="navlink">
    <div class="container">
        <ul class="clearfix">
            <li ><a href="/">首页</a></li>

            <li ><a href="/list/remen" title="微信热门">热门</a></li>

            <li ><a href="/list/tuijian" title="微信推荐">推荐</a></li>

            <li ><a href="/list/gaoxiao" title="微信搞笑">搞笑</a></li>

            <li ><a href="/list/jiankang" title="微信健康">健康</a></li>

            <li ><a href="/list/liangxing" title="微信两性">两性</a></li>

            <li ><a href="/list/bagua" title="微信八卦">八卦</a></li>

            <li ><a href="/list/shenghuo" title="微信生活">生活</a></li>

            <li ><a href="/list/caijing" title="微信财经">财经</a></li>

            <li ><a href="/list/qiche" title="微信汽车">汽车</a></li>

            <li ><a href="/list/keji" title="微信科技">科技</a></li>

            <li ><a href="/list/nvxing" title="微信女性">女性</a></li>

            <li ><a href="/list/lama" title="微信辣妈">辣妈</a></li>

            <li ><a href="/list/lizhi" title="微信励志">励志</a></li>

            <li ><a href="/list/chaoliu" title="微信潮流">潮流</a></li>

            <li ><a href="/list/zhichang" title="微信职场">职场</a></li>

            <li ><a href="/list/meishi" title="微信美食">美食</a></li>

            <li ><a href="/list/lishi" title="微信历史">历史</a></li>

            <li ><a href="/list/jiaoyu" title="微信教育">教育</a></li>

            <li ><a href="/list/xingzuo" title="微信星座">星座</a></li>

            <li ><a href="/list/tiyu" title="微信体育">体育</a></li>

            <li ><a href="/update" title="最近更新">最新</a></li>
            <li ><a href="/wxlist" title="微信公共号">微信</a></li>
            <li><a href="http://book.qixi.us" title="免费小说">小说</a></li>
        </ul>
    </div>
</div>



<div class="wrap">

    <div class="container clearfix">
        <div class="page-container">
            <div class="page-main">

                <div class="crumbs">
                    <a href="/">首页</a>
                    <i>&rsaquo;</i>
                    <span>关于我们</span>
                </div>

                <div class="about">


                    <div class="panel">
                        <div class="panel-title">
                            <h1 class="title"><i class="icon"></i>关于我们</h1>
                        </div>
                        <div class="panel-content">
                            <div class="content">
                                <p>
                                    七夕阅读网是基于数据挖掘的智能推荐引擎，聚合微信公共号文章门户，微信热点资讯站，微信今日头条，及微信小说等阅读平台。
                                </p>
                            </div>
                        </div>
                    </div>












                </div>


            </div>
            <div class="page-side">


                <div class="panel">
                    <div class="panel-title">
                        <h3 class="title"><i class="icon"></i>快捷导航</h3>
                    </div>
                    <div class="panel-content">
                        <ul class="catalog-list">
                            <li><a href="/about/aboutus">关于我们</a></li>
                            <li><a href="/about/copyright">版权声明</a></li>
                            <li><a href="/about/contact">联系我们</a></li>
                            <li><a href="/about/biz">商务合作</a></li>
                        </ul>
                    </div>
                </div>


            </div>
        </div>
    </div>

</div>



<div class="footer">
    <div class="container">
        <p>
            <a href="/about/aboutus">关于我们</a>
            <span>|</span>
            <a href="/about/copyright">版权声明</a>
            <span>|</span>
            <a href="/about/contact">联系我们</a>
            <span>|</span>
            <a href="/about/biz">商务合作</a>
        </p>

        <p>&copy 七夕阅读网(www.qixi.us) <a href="#">粤ICP备07000266号-8</a></p>
        <p class="text-danger">本站为第三方微信站，本站与腾讯微信、微信公众平台无任何关联，非腾讯微信官方网站。</p>
        <p class="copyright text-left">
            版权声明：本站收录微信公众号和微信文章内容全部来自于网络，仅供个人学习、研究或者欣赏使用，版权归原作者所有。<br>
            如果您发现网站上有侵犯您的知识产权的内容，请与我们联系，我们会及时修改或删除。
        </p>
    </div>
</div>













<div class="stat">
    <script src="https://s19.cnzz.com/z_stat.php?id=1262237922&web_id=1262237922" language="JavaScript"></script>
</div>

</body>

</html>`