package main

import (
	"log"
	"net/http"

	"github.com/Deichindianer/poe-item-bot/internal/itemservice"
)

func main() {
	is := itemservice.NewItemService("SSF Ritual HC")
	go is.Init()
	log.Fatal(http.ListenAndServe(":8080", is))
}
