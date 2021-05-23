package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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
		})

		// ファイルに取得したデータを保存する

		file, err := os.Create("./json/players.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for _, player := range players {
			p, _ := json.Marshal(player)
			_, err = file.Write(p)
			if err != nil {
				log.Fatal(err)
			}
		}

		return c.JSON(http.StatusOK, players)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
