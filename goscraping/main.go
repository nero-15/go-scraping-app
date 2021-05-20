package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}"` + "\n",
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {

		url := "https://www.f-marinos.com/"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(byteArray))

		return c.HTML(http.StatusOK, string(byteArray))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
