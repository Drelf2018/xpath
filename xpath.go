package xpath

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

type Any struct {
	Type  unsafe.Pointer
	Value unsafe.Pointer
}

func TypePtr(in any) uintptr {
	return uintptr((*Any)(unsafe.Pointer(&in)).Type)
}

// ValuePtr can obtain the uintptr of a type from its reflect.Type
//
//	TypePtr(something{}) is equal to ValuePtr(reflect.TypeOf(something{}))
func ValuePtr(in any) uintptr {
	return uintptr((*Any)(unsafe.Pointer(&in)).Value)
}

type Expr struct {
	Index int
	*xpath.Expr
}

func GetExpr(elem reflect.Type) ([]*Expr, error) {
	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xpath.GetExpr: non-struct %v", elem)
	}
	r := make([]*Expr, 0, elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("xpath")
		if tag == "" {
			continue
		}

		expr, err := xpath.Compile(tag)
		if err != nil {
			return nil, err
		}

		r = append(r, &Expr{i, expr})
	}
	return r, nil
}

var exprCache sync.Map // map[uintptr][]*Expr

func LoadExpr(v any) ([]*Expr, error) {
	ptr := TypePtr(v)
	if expr, ok := exprCache.Load(ptr); ok {
		return expr.([]*Expr), nil
	}
	elem := reflect.TypeOf(v)
	if elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	expr, err := GetExpr(elem)
	if err != nil {
		return nil, err
	}
	exprCache.Store(ptr, expr)
	return expr, nil
}

func unmarshalValue(node *html.Node, v reflect.Value, expr *xpath.Expr) (err error) {
	if v.Kind() != reflect.Slice {
		return unmarshal(htmlquery.QuerySelector(node, expr), v)
	}

	nodes := htmlquery.QuerySelectorAll(node, expr)
	if v.Len() < len(nodes) {
		v.Set(reflect.MakeSlice(v.Type(), len(nodes), len(nodes)))
	}

	isPtr := v.Type().Elem().Kind() == reflect.Pointer
	for idx := range nodes {
		item := v.Index(idx)
		if isPtr {
			item.Set(reflect.New(item.Type().Elem()))
		}
		err = unmarshal(nodes[idx], item)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshal(node *html.Node, v reflect.Value) error {
	if node == nil {
		return nil
	}

	if v.Kind() == reflect.String {
		if node.Type == html.ElementNode {
			if node.LastChild != nil {
				v.SetString(node.LastChild.Data)
			}
		} else {
			v.SetString(node.Data)
		}
		return nil
	}

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("xpath.unmarshal: non-struct %v", v.Type())
	}

	expr, err := LoadExpr(v.Interface())
	if err != nil {
		return err
	}

	for _, e := range expr {
		err = unmarshalValue(node, v.Field(e.Index), e.Expr)
		if err != nil {
			return err
		}
	}

	return nil
}

type XPath interface {
	XPath() string
}

func UnmarshalNode(node *html.Node, v any) error {
	if v, ok := v.(XPath); ok {
		expr, err := xpath.Compile(v.XPath())
		if err != nil {
			return err
		}
		return unmarshalValue(node, reflect.ValueOf(v).Elem(), expr)
	}
	return unmarshal(node, reflect.ValueOf(v))
}

func UnmarshalNodeWith[V any](node *html.Node) (v V, err error) {
	err = UnmarshalNode(node, &v)
	return
}

func UnmarshalReader(r io.Reader, v any) error {
	node, err := html.Parse(r)
	if err != nil {
		return err
	}
	return UnmarshalNode(node, v)
}

func UnmarshalReaderWith[V any](r io.Reader) (v V, err error) {
	err = UnmarshalReader(r, &v)
	return
}

func UnmarshalText(text string, v any) error {
	return UnmarshalReader(strings.NewReader(text), v)
}

func UnmarshalTextWith[V any](text string) (v V, err error) {
	err = UnmarshalText(text, &v)
	return
}

func Unmarshal(b []byte, v any) error {
	return UnmarshalReader(bytes.NewReader(b), v)
}

func UnmarshalWith[V any](b []byte) (v V, err error) {
	err = Unmarshal(b, &v)
	return
}

func LoadURL(url string, v any) error {
	node, err := htmlquery.LoadURL(url)
	if err != nil {
		return err
	}
	return UnmarshalNode(node, v)
}

func LoadURLWith[V any](url string) (v V, err error) {
	err = LoadURL(url, &v)
	return
}

func LoadDoc(path string, v any) error {
	node, err := htmlquery.LoadDoc(path)
	if err != nil {
		return err
	}
	return UnmarshalNode(node, v)
}

func LoadDocWith[V any](path string) (v V, err error) {
	err = LoadDoc(path, &v)
	return
}
