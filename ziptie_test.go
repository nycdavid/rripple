package ziptie

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/labstack/echo.v3"
)

type PostsCtrl struct {
	Namespace string
	Index     interface{} `path:"" method:"GET"`
	Show      interface{} `path:"/:id" method:"GET"`
}

func (ctrl *PostsCtrl) IndexFunc(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Index")
}

func (ctrl *PostsCtrl) ShowFunc(ctx echo.Context) error {
	res := fmt.Sprintf("Showing post %s", ctx.Param("id"))
	return ctx.String(http.StatusOK, res)
}

func TestFastenWithOneMethod(t *testing.T) {
	e := echo.New()
	Fasten(&PostsCtrl{Namespace: "/posts"}, e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Error("Status Code is not 200")
	}
}

func TestFastenWithASecondMethod(t *testing.T) {
	e := echo.New()
	Fasten(&PostsCtrl{Namespace: "/posts"}, e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	e.ServeHTTP(rec, req)
	byt := rec.Body.Bytes()

	if rec.Code != 200 {
		t.Error("Status Code is not 200")
	}
	if string(byt) != "Showing post 1" {
		t.Error("Unexpected response body")
	}
}

type MixedCtrl struct {
	Namespace string
	Config    map[string]interface{}
	Index     interface{} `path:"" method:"GET"`
}

func (mc *MixedCtrl) IndexFunc(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Index")
}

func TestHandlingOfNonMethodFieldsInStruct(t *testing.T) {
	e := echo.New()
	Fasten(&MixedCtrl{Namespace: "/mixed"}, e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/mixed", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Error("Status Code is not 200")
	}
}

type MissingHandlerFuncCtrl struct {
	Index interface{} `path:"" method:"GET"`
}

func TestMissingHandlerFunc(t *testing.T) {
	t.Skip("TODO")
}

type NamespacedCtrl struct {
	Namespace string
	Foo       interface{} `path:"/" method:"GET"`
}

func (nc *NamespacedCtrl) FooFunc(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "foo")
}

func TestRootCtrl(t *testing.T) {
	e := echo.New()
	ctrl := &NamespacedCtrl{Namespace: ""}
	Fasten(ctrl, e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(echo.GET, "/", nil)
	e.ServeHTTP(rec, req)
	byt := rec.Body.Bytes()

	if rec.Code != 200 {
		t.Error(fmt.Sprintf("Expecting %d, got %d", 200, rec.Code))
	}
	if string(byt) != "foo" {
		t.Error("Unexpected response body")
	}
}
