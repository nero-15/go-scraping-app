package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Player struct {
	Url      string `json:"url"`
	Number   string `json:"number"`
	Position string `json:"position"`
	Img      string `json:"img"`
	NameJp   string `json:"nameJp"`
	NameEn   string `json:"nameEn"`
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

		players := []Player{}
		doc.Find("a.card-player").Each(func(i int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			number := s.Find(".card-player-number").Text()
			position := s.Find(".card-player-position").Text()
			img, _ := s.Find("img").Attr("src")
			nameJp := s.Find(".card-player-name-jp").Text()
			nameEn := s.Find(".card-player-name-en").Text()

			player := Player{
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
		jsonData, _ := json.Marshal(players)
		fmt.Println(jsonData)

		return c.JSON(http.StatusOK, string(jsonData))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
