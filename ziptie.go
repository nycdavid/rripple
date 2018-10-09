package ziptie

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo"
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
		if field.Name == "Config" || field.Name == "Namespace" || field.Name == "DB" {
			continue
		}
		path, method := extractPathAndMethod(field)
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

func extractPathAndMethod(field reflect.StructField) (string, string) {
	return field.Tag.Get("path"), field.Tag.Get("method")
}
