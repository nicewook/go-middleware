package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo"
)

type Response struct {
	Code    int         `json:"code`
	Message string      `json:"message"`
	User    interface{} `json:"user`
}

type LogFilter struct {
	Level   string `json:"level"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserID  string `json:"UserID"`
}

type User struct {
	UserID  string `json:"userid" logging:"true"`
	Email   string `json:"email" logging:"true"`
	Age     int    `json:"age" logging:"true"`
	Name    string `json:"name"`
	Address string `json:"address"`
	GoodsID string `json:"goodsid" logging:"goodsid"`
}

func (user *User) Set(c echo.Context) {
	c.Set("user", user)
}

func logMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		var clog []byte
		start := time.Now()
		if err = next(c); err != nil {
			return err
		}
		diff := time.Since(start)
		req := c.Request()
		resp := c.Response()
		user := c.Get("user")
		if user != nil {
			t := reflect.ValueOf(user)
			e := t.Elem()
			m := make(map[string]interface{})
			for i := 0; i < e.NumField(); i++ {
				mValue := e.Field(i)
				mType := e.Type().Field(i)
				if _, ok := mType.Tag.Lookup("logging"); ok {
					m[mType.Name] = mValue.Interface()
				}
			}
			clog, _ = json.Marshal(m)
		}
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
	u := User{
		UserID:  "10a293d",
		Email:   "yundream@gmail.com",
		Age:     38,
		Name:    "yundream",
		Address: "seoul",
		GoodsID: "3849-14281",
	}
	u.Set(c)
	return c.JSON(http.StatusOK, Response{101, "pong", u})
}

func main() {
	e := echo.New()
	e.Use(logMiddleWare)
	e.GET("/ping", Ping)
	e.Logger.Fatal(e.Start(":1234"))
}
