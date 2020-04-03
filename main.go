package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//updateMill()
	checkMill()
}

func checkMill() {
	sel := getMill()

	b, _ := ioutil.ReadFile("mill.txt")

	state := string(b)

	if sel != state {
		log.Printf("The page has changed!\n")
	} else {
		log.Printf("Still the same\n")
	}
}

func updateMill() {
	sel := getMill()

	_ = ioutil.WriteFile("mill.txt", []byte(sel), 0644)
	log.Printf("updated file\n")
}

func getMill() string {
	res, err := http.Get("https://www.shipton-mill.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got the page\n")

	var strong string

	doc.Find(".well").Each(func(i int, s *goquery.Selection) {
		strong = s.Find("p").Text()
	})

	return strong
}
