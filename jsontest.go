package main

import (
	"fmt"
	"reflect"

	"github.com/bitly/go-simplejson"
)

func main() {
	err := t1()
	if err != nil {
		fmt.Println(err)
	}
}

//golang 加载jison而无需定义一堆结构体
func t1() error {
	jsonStr := `{"aa":[{"a":333},3,4],"bb":"111","cc":{"dd":333}}`
	json, err := simplejson.NewJson([]byte(jsonStr))
	if err != nil {
		return err
	}
	var nodes = make(map[string]interface{})
	nodes, _ = json.Map()

	fmt.Println(jsonStr)
	fmt.Println(nodes)
	fmt.Println(nodes["aa"])
	fmt.Println(reflect.TypeOf(nodes["aa"]).Kind())
	switch reflect.TypeOf(nodes["aa"]).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(nodes["aa"])
		fmt.Println(s.Index(0))
		tt, ok := s.Index(0).Interface().(map[string]interface{})
		if ok {
			fmt.Println(tt["a"])
		}
		fmt.Println(tt["a"])
		fmt.Println(reflect.TypeOf(s.Index(0)).Kind())
	}
	return nil
}
