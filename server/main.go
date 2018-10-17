package main

import (
	"log"
	"net/http"

	app "github.com/mrmod/magicland/src"
)

func main() {
	http.HandleFunc("/", app.WebhookHandler)
	log.Println("Running on port 8000")
	http.ListenAndServe(":8000", nil)
}
