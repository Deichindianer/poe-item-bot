package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	_ "net/http/pprof"

	"github.com/Deichindianer/poe-item-bot/internal/itemservice"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	is := itemservice.NewItemService("SSF Ritual HC", 45, 0)
	go is.Init()
	log.Fatal(http.ListenAndServe(":8080", is))
}
