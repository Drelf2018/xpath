package xpath

import (
	"errors"
	"reflect"

	"github.com/Drelf2018/reflectMap"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

func unmarshal(node *html.Node, v reflect.Value, datas []reflectMap.Data[*xpath.Expr]) {
	if node == nil {
		return
	}

	if v.Kind() == reflect.String {
		if node.Type == html.ElementNode {
			v.SetString(node.LastChild.Data)
		} else {
			v.SetString(node.Data)
		}
		return
	}

	for _, data := range datas {
		field := v.Field(data.Index)

		if field.Kind() != reflect.Slice {
			unmarshal(htmlquery.QuerySelector(node, data.V), field, data.Fields)
			continue
		}

		nodes := htmlquery.QuerySelectorAll(node, data.V)
		nodesLen := len(nodes)
		if field.Len() < nodesLen {
			field.Set(reflect.MakeSlice(field.Type(), nodesLen, nodesLen))
		}
		for j := range nodes {
			unmarshal(nodes[j], field.Index(j), data.Fields)
		}
	}
}

func convert(fieldType reflect.Type) reflect.Type {
	if fieldType.Kind() == reflect.Slice {
		return fieldType.Elem()
	}
	return fieldType
}

var GetData = reflectMap.NewTagParser("xpath", xpath.MustCompile).SetConverter(convert).Get
var ErrNotPtr = errors.New("xpath: parameter v must be a pointer to struct")

func Unmarshal(data *html.Node, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return ErrNotPtr
	}
	unmarshal(data, val.Elem(), GetData(v))
	return nil
}

func LoadURL(url string, v any) error {
	node, err := htmlquery.LoadURL(url)
	if err != nil {
		return err
	}
	return Unmarshal(node, v)
}

func MustUnmarshal[V any](data *html.Node) *V {
	v := new(V)
	Unmarshal(data, v)
	return v
}

func MustLoadURL[V any](url string) *V {
	v := new(V)
	LoadURL(url, v)
	return v
}
