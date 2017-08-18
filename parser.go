package scrago

import (
	"github.com/PuerkitoBio/goquery"
	"reflect"
	"fmt"
	"regexp"
	"strings"
	"strconv"
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
	//log.Printf("%#v kind is %v | %v", v, refValue.Kind(), reflect.Ptr)
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
		//log.Printf("===== node selector : %v, func: %v, params: %v", cssParser.Selector, cssParser.FuncName, cssParser.FuncParams)
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
			if callReturns[0].Type() != fieldType.Type{
				return fmt.Errorf("method %v return value of type %v is not assignable to type %v", cssParser.FuncName, callReturns[0].Type(), fieldType.Type)
			}
			if len(callReturns) > 1{
				if err, ok := callReturns[len(callReturns) - 1].Interface().(error); ok{
					if err != nil {
						return fmt.Errorf("method %v return error: %v", cssParser.FuncName, err)
					}
				}
			}
			fieldValue.Set(callReturns[0])
			continue
		}

		//set value
		kind := fieldType.Type.Kind()
		switch {
		//Bool
		case kind == reflect.Bool:
			value, err := strconv.ParseBool(nodeValue)
			if err != nil {
				return fmt.Errorf("field %v convert value %v to %v error: %v", fieldType, nodeValue, kind, err)
			}
			fieldValue.SetBool(value)
			//Int
			//Int8
			//Int16
			//Int32
			//Int64
		case kind >= reflect.Int && kind <= reflect.Int64:
			value, err := strconv.ParseInt(nodeValue, 10, int(fieldValue.Type().Size()*8))
			if err != nil {
				return fmt.Errorf("field %v convert value %v to %v error: %v", fieldType, nodeValue, kind, err)
			}
			fieldValue.SetInt(value)
			//Uint
			//Uint8
			//Uint16
			//Uint32
			//Uint64
			//Uintptr
		case kind >= reflect.Uint && kind <= reflect.Uintptr:
			value, err := strconv.ParseUint(nodeValue, 10, int(fieldValue.Type().Size()*8))
			if err != nil {
				return fmt.Errorf("field %v convert value %v to %v error: %v", fieldType, nodeValue, kind, err)
			}
			fieldValue.SetUint(value)
			//Float32
			//Float64
		case kind == reflect.Float32 || kind == reflect.Float64:
			value, err := strconv.ParseFloat(nodeValue, 64)
			if err != nil {
				return fmt.Errorf("field %v convert value %v to %v error: %v", fieldType, nodeValue, kind, err)
			}
			fieldValue.SetFloat(value)
			//Interface
		case kind == reflect.Interface:
			fieldValue.Set(reflect.ValueOf(nodeValue))
			//Map
			//Ptr
		case kind == reflect.Ptr:
			subModel := reflect.New(fieldType.Type.Elem())
			fieldValue.Set(subModel)
			err = ParserField(subModel.Interface(), node)
			if err != nil {
				return fmt.Errorf("%#v parser error: %v", subModel, err)
			}
			//Slice
		case kind == reflect.Slice:
			slicetyp := fieldValue.Type()
			itemtyp := slicetyp.Elem()
			itemkind := itemtyp.Kind()
			slice := reflect.MakeSlice(slicetyp, node.Size(), node.Size())
			node.EachWithBreak(func(i int, subNode *goquery.Selection) bool {
				//outhtml, _ := goquery.OuterHtml(subNode)
				//log.Printf("%v => %v", i, outhtml)
				tmp := reflect.New(itemtyp).Elem()
				switch {
				case itemkind == reflect.String:
					tmp.SetString(subNode.Text())
				case itemkind == reflect.Struct:
					err = ParserField(tmp.Addr().Interface(), subNode)
					if err != nil {
						err = fmt.Errorf("%#v parser error: %v", tmp, err)
						return false
					}
				case itemkind == reflect.Ptr && tmp.Type().Elem().Kind() == reflect.String:
					tmpStr := subNode.Text()
					tmp.Set(reflect.ValueOf(&tmpStr))
				case itemkind == reflect.Ptr && tmp.Type().Elem().Kind() == reflect.Struct:
					tmp = reflect.New(itemtyp.Elem())
					err = ParserField(tmp.Interface(), subNode)
					if err != nil {
						err = fmt.Errorf("%#v parser error: %v", tmp, err)
						return false
					}
				default:
					err = fmt.Errorf("slice not support item type %v, kind %v", itemtyp, itemkind)
					return false
				}
				slice.Index(i).Set(tmp)
				return true
			})
			if err != nil {
				return err
			}
			fieldValue.Set(slice)
		case kind == reflect.String:
			fieldValue.SetString(nodeValue)
		case kind == reflect.Struct:
			subModel := reflect.New(fieldType.Type)
			err = ParserField(subModel.Interface(), node)
			if err != nil {
				return fmt.Errorf("%#v parser error: %v", subModel, err)
			}
			fieldValue.Set(subModel.Elem())
			//UnsafePointer
			//Complex64
			//Complex128
			//Array
			//Chan
			//Func
		default:
			return fmt.Errorf("field %v not support type %v", fieldType, kind)
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
