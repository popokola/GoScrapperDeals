package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Deal struct {
	Name     string `json:"name"`
	Link     string `json:"link"`
	Temp     string `json:"temp"`
	EndPromo string `json:"end_promo"`
}

func writeJSONFile(deals []Deal) {
	json, err := json.MarshalIndent(deals, "", "  ")
	if err != nil {
		log.Fatal(err)
		return
	}

	_ = ioutil.WriteFile("deals.json", json, 0644)

}

func main() {
	deals := []Deal{}
	numPage := 5

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		//set the user agent
		//r.Headers.Set("User-Agent", "Put Your User Agent Here")
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response code", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error", err.Error())
	})

	c.OnHTML(".thread", func(e *colly.HTMLElement) {
		div := e.DOM
		temp := div.Find(".cept-vote-temp").Text()
		temp = strings.TrimSpace(temp)
		if temp == "" {
			return
		} else {
			temp = temp[:len(temp)-2]
		}

		//keep only above 500
		if temp < "500" {
			return
		}

		link := div.Find("a").AttrOr("href", "no link")

		name := div.Find(".thread-title").Text()

		endPromo := div.Find(".hide--toW3").Text()
		split := strings.Split(endPromo, "min")
		if len(split) > 1 {
			endPromo = split[0]
		}

		d := Deal{Name: name, Temp: temp, Link: link, EndPromo: endPromo}
		deals = append(deals, d)
	})

	for i := 1; i <= numPage; i++ {
		c.Visit("https://dealabs.com/?page=" + strconv.Itoa(i))
	}

	writeJSONFile(deals)

}
