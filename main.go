package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	// for checking the mills
	//	work()
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func init() {
	if ok := os.Getenv("FREQUENCY"); ok == "" {
		log.Fatalln("FREQUENCY is not defined")
	}
}
