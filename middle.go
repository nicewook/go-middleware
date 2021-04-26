package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type Response struct {
	Code    int    `json:"code`
	Message string `json:"message"`
}

type LogFilter struct {
	Level   string `json:"level"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserID  string `json:"UserID"`
}

func simpleMiddelWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("-- It is the simple and first middleware")
		err := next(c)
		return err
	}

}

func logMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		start := time.Now()
		if err = next(c); err != nil {
			return err
		}
		diff := time.Since(start)
		req := c.Request()
		resp := c.Response()
		clog, _ := json.Marshal(c.Get("clog"))
		fmt.Println(
			start.Format(time.RFC3339),
			diff,
			c.RealIP(),
			req.Method,
			req.Host,
			req.Header.Get(echo.HeaderContentLength),
			resp.Size,
			string(clog),
		)
		return err
	}
}

func Ping(c echo.Context) error {
	time.Sleep(20 * time.Millisecond)
	c.Set("clog", LogFilter{"INFO", 101, "pong", "hsjeong@gmail.com"})
	return c.JSON(http.StatusOK, Response{101, "pong"})
}

func main() {
	e := echo.New()
	e.Use(simpleMiddelWare)
	e.Use(logMiddleWare)
	e.GET("/ping", Ping)
	e.Logger.Fatal(e.Start(":1234"))
}
