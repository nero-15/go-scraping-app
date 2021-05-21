package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type player struct {
	url      string
	number   string
	position string
	img      string
	nameJp   string
	nameEn   string
}

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}"` + "\n",
	}))
	e.Use(middleware.Recover())

	e.GET("/player", func(c echo.Context) error { //選手・スタッフページから選手のデータを取得する

		url := "https://www.f-marinos.com/team/player"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		players := []player{}
		doc.Find("a.card-player").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			number := s.Find(".card-player-number").Text()
			position := s.Find(".card-player-position").Text()
			img, _ := s.Find("img").Attr("src")
			nameJp := s.Find(".card-player-name-jp").Text()
			nameEn := s.Find(".card-player-name-en").Text()

			player := player{
				url,
				number,
				position,
				img,
				nameJp,
				nameEn,
			}

			players = append(players, player)

			fmt.Printf("Review %d: %s - %s - %s - %s - %s - %s\n", i, url, number, position, img, nameJp, nameEn)
		})
		fmt.Println(players)
		return c.String(http.StatusOK, "player")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
