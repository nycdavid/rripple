package rripple

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/labstack/echo.v3"
)

type Ctrl interface{}

func Group(ctrl Ctrl, e *echo.Echo) {
	vof := reflect.ValueOf(ctrl)
	elem := vof.Elem()
	typ := vof.Type().String()
	route := convertRouteFromType(typ)

	numOfRoutes := elem.Type().NumField()
	for i := 0; i < numOfRoutes; i++ {
		field := elem.Type().Field(i)
		path := field.Tag.Get("path")
		method := field.Tag.Get("method")
		handler := vof.MethodByName(fmt.Sprintf("%sFunc", field.Name))
		echMd := reflect.ValueOf(e).MethodByName(method)
		args := []reflect.Value{
			reflect.ValueOf(fmt.Sprintf("%s%s", route, path)),
			handler,
		}
		echMd.Call(args)
	}
}

func convertRouteFromType(ctrlType string) string {
	ctrlName := strings.Split(ctrlType, ".")[1]
	rsrcName := strings.Replace(ctrlName, "Ctrl", "", 1)
	return fmt.Sprintf("/%s", strings.ToLower(rsrcName))
}
