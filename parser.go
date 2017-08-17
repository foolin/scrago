package scrago

import (
	"github.com/PuerkitoBio/goquery"
	"reflect"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"log"
)
//::text()
//::value()
//::html()
//::attr(xxx)
var rxFunc = regexp.MustCompile("^\\s*([a-zA-Z]+)\\s*\\(([^\\)]*)\\){0,1}\\s*$")
const (
	parserTagName = "css"
	parserSplitSep = "::"
)

type CssParser struct {
	Selector string
	FuncName string
	FuncParams []string
}


func ParserField(v interface{}, selection *goquery.Selection) (err error) {
	refType := reflect.TypeOf(v)
	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Ptr{
		return fmt.Errorf("%v is non-pointer?",  refType)
	}
	if refValue.IsNil(){
		return fmt.Errorf("%v is nil?", refType)
	}
	refTypeElem := refType.Elem()
	refValueElem := refValue.Elem()
	for i := 0; i < refValueElem.NumField(); i++{
		fieldType := refTypeElem.Field(i)
		fieldValue := refValueElem.Field(i)
		tagValue := fieldType.Tag.Get(parserTagName)
		cssParser := newCssParser(tagValue)
		log.Printf("===== node selector : %v, func: %v, params: %v", cssParser.Selector, cssParser.FuncName, cssParser.FuncParams)
		node := selection
		if cssParser.Selector != ""{
			node = selection.Find(cssParser.Selector)
		}
		nodeValue := ""
		switch strings.ToLower(cssParser.FuncName) {
		case "text":
			nodeValue = node.Text()
		case "html":
			nodeValue, _ = node.Html()
		case "outerhtml":
			nodeValue, _ = goquery.OuterHtml(node)
		case "value":
			nodeValue = node.AttrOr("value", "")
		case "attr":
			if len(cssParser.FuncParams) > 0 && cssParser.FuncParams[0] != ""{
				nodeValue = node.AttrOr(cssParser.FuncParams[0], "")
			}
		case "":
			nodeValue = node.Text()
		default:
			callMethod := refValueElem.MethodByName(cssParser.FuncName)
			if !callMethod.IsValid(){
				callMethod = refValue.MethodByName(cssParser.FuncName)
			}
			if !callMethod.IsValid(){
				return fmt.Errorf("method %v not found!", cssParser.FuncName)
			}
			callParams := make([]reflect.Value, 0)
			callParams = append(callParams, reflect.ValueOf(node))
			callReturns := callMethod.Call(callParams)
			if len(callReturns) <= 0{
				return fmt.Errorf("method %v not return any value", cssParser.FuncName)
			}
			fieldValue.Set(callReturns[0])
			continue
		}

		switch fieldType.Type.Kind() {
			//Bool
		case reflect.Bool:
			value, _ := strconv.ParseBool(nodeValue)
			fieldValue.SetBool(value)
			//Int
			//Int8
			//Int16
			//Int32
			//Int64
			//Uint
			//Uint8
			//Uint16
			//Uint32
			//Uint64
			//Uintptr
			//Float32
			//Float64
			//Complex64
			//Complex128
			//Array
			//Chan
			//Func
			//Interface
			//Map
			//Ptr
			//Slice

		case reflect.String:
			fieldValue.SetString(nodeValue)
		case reflect.Struct:
			subModel := reflect.New(fieldType.Type)
			err = ParserField(subModel, node)
			if err != nil {
				return fmt.Errorf("%v parser error: %v", fieldType, err)
			}
			//UnsafePointer
		}
	}
	return nil
}

func newCssParser(tagValue string) *CssParser {
	cssParser := &CssParser{}
	if tagValue == ""{
		return cssParser
	}
	selectors := strings.Split(tagValue, parserSplitSep)
	funcValue := ""
	for i := 0; i < len(selectors); i++{
		switch i {
		case 0:
			cssParser.Selector = strings.TrimSpace(selectors[i])
		case 1:
			funcValue = selectors[i]
		}
	}
	matchs := rxFunc.FindStringSubmatch(funcValue)
	if len(matchs) < 3{
		return cssParser
	}
	cssParser.FuncName = strings.TrimSpace(matchs[1])
	cssParser.FuncParams = strings.Split(matchs[2], ",")
	return cssParser
}