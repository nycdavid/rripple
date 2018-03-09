package ziptie

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/labstack/echo.v3"
)

type Ctrl interface{}

func Fasten(ctrl Ctrl, e *echo.Echo) {
	vof := reflect.ValueOf(ctrl)
	elem := vof.Elem()
	namespace := elem.FieldByName("Namespace").String()
	route := convertRouteFromType(namespace)

	numOfRoutes := elem.Type().NumField()
	for i := 0; i < numOfRoutes; i++ {
		field := elem.Type().Field(i)
		if field.Name == "Config" || field.Name == "Namespace" {
			continue
		}
		path := field.Tag.Get("path")
		method := field.Tag.Get("method")
		handler := vof.MethodByName(fmt.Sprintf("%sFunc", field.Name))
		echMd := reflect.ValueOf(e).MethodByName(method)
		fullPath := fmt.Sprintf("%s%s", route, path)
		args := []reflect.Value{
			reflect.ValueOf(fullPath),
			handler,
		}
		echMd.Call(args)
	}
}

func convertRouteFromType(namespace string) string {
	return fmt.Sprintf("%s", strings.ToLower(namespace))
}
