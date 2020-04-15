package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	SHIPTON  = "https://www.shipton-mill.com"
	MATTHEWS = "https://www.fwpmatthews.co.uk/"
)

func main() {
	updateMill(SHIPTON)
	updateMill(MATTHEWS)
	worker()
}

func worker() {
	log.Printf("Dusting flour...\n")
	freq, _ := strconv.Atoi(os.Getenv("FREQUENCY"))
	f := time.Duration(freq)
	uptimeTicker := time.NewTicker(f * time.Second)

	for {
		select {
		case <-uptimeTicker.C:
			checkMill(SHIPTON)
			checkMill(MATTHEWS)
		}
	}
}

func checkMill(mill string) {
	sel := getMill(mill)

	if mill == SHIPTON {
		b, _ := ioutil.ReadFile("shiptonmill.txt")

		state := string(b)

		log.Printf("SHIPTON MILL\n")
		if sel != state {
			log.Printf("The page has changed!\n")
			err := notify(SHIPTON, state, sel)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatalf("This is a hard exit to not overflood me with notifications\n")
		} else {
			log.Printf("Still the same\n")
		}
	}
	if mill == MATTHEWS {
		b, _ := ioutil.ReadFile("matthewsmill.txt")

		state := string(b)

		log.Printf("MATTHEWS MILL\n")
		if sel != state {
			log.Printf("The page has changed!\n")
			err := notify(MATTHEWS, state, sel)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatalf("This is a hard exit to not overflood me with notifications\n")
		} else {
			log.Printf("Still the same\n")
		}
	}
}

func updateMill(mill string) {
	sel := getMill(mill)

	if mill == SHIPTON {
		_ = ioutil.WriteFile("shiptonmill.txt", []byte(sel), 0644)
		log.Printf("updated shipton mill file\n")
	}
	if mill == MATTHEWS {
		_ = ioutil.WriteFile("matthewsmill.txt", []byte(sel), 0644)
		log.Printf("updated matthews mill file\n")
	}
}

func getMill(mill string) string {
	res, err := http.Get(mill)
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

	var result string

	if mill == SHIPTON {
		doc.Find(".well").Each(func(i int, s *goquery.Selection) {
			result = s.Find("p").Text()
		})
	}
	if mill == MATTHEWS {
		doc.Find(".storeclosing_popup").Each(func(i int, s *goquery.Selection) {
			result = s.Find("div").Text()
		})
	}
	return result
}

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func notify(mill, state, content string) error {
	from := os.Getenv("NOTIFICATION_EMAIL_SEND")
	password := os.Getenv("NOTIFICATION_EMAIL_SEND_PASSWORD")

	receiver := []string{os.Getenv("NOTIFICATION_EMAIL_RECEIVER")}

	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	var message []byte

	if mill == SHIPTON {
		message = []byte(fmt.Sprintf("shipton mill has changed something!\n\npreviously: %v\nnow: %v\n", state, content))
	}

	if mill == MATTHEWS {
		message = []byte(fmt.Sprintf("matthews mill has changed something!\n\n%vpreviously: %v\nnow: %v\n", state, content))
	}

	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	return smtp.SendMail(smtpServer.Address(), auth, from, receiver, message)
}
