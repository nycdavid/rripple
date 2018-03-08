package ziptie

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/labstack/echo.v3"
)

type PostsCtrl struct {
	Index interface{} `path:"" method:"GET"`
	Show  interface{} `path:"/:id" method:"GET"`
}

func (ctrl *PostsCtrl) IndexFunc(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Index")
}

func (ctrl *PostsCtrl) ShowFunc(ctx echo.Context) error {
	res := fmt.Sprintf("Showing post %s", ctx.Param("id"))
	return ctx.String(http.StatusOK, res)
}

func TestRrippleGroupWithOneMethod(t *testing.T) {
	e := echo.New()
	Fasten(&PostsCtrl{}, e)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Error("Status Code is not 200")
	}
}

func TestRrippleGroupWithASecondMethod(t *testing.T) {
	e := echo.New()
	Fasten(&PostsCtrl{}, e)

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
