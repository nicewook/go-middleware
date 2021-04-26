package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type Response struct {
	Code    int    `json:"code`
	Message string `json:"message"`
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{101, "pong"})
}

func main() {
	e := echo.New()
	e.GET("/ping", Ping)
	e.Logger.Fatal(e.Start(":1234"))
}
