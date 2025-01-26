package main

import (
	"github.com/labstack/echo/v4"
	"github.com/nayuneko/cgi2echo"
	"net/http"
)

func helloHandler(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	return c.String(http.StatusOK, "Hello World")
}

func main() {
	c := cgi2echo.New()
	c.Echo().GET("/api/hello", helloHandler)
	if err := c.Serve(); err != nil {
		panic(err)
	}
}
