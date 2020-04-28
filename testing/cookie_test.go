package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestCookiesBasic(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://example.com"))

	cookieName, cookieValue := "my_cookie_name", "my_cookie_value"

	// Test Set A Cookie.
	t1 := e.GET(fmt.Sprintf("/cookies/%s/%s", cookieName, cookieValue)).
		Expect().Status(httptest.StatusOK)
	// Validate cookie's existence, it should be available now.
	t1.Cookie(cookieName).Value().Equal(cookieValue)
	t1.Body().Contains(cookieValue)

	path := fmt.Sprintf("/cookies/%s", cookieName)

	// Test Retrieve A Cookie.
	t2 := e.GET(path).Expect().Status(httptest.StatusOK)
	t2.Body().Equal(cookieValue)

	// Test Remove A Cookie.
	t3 := e.DELETE(path).Expect().Status(httptest.StatusOK)
	t3.Body().Contains(cookieName)

	t4 := e.GET(path).Expect().Status(httptest.StatusOK)
	t4.Cookies().Empty()
	t4.Body().Empty()
}
