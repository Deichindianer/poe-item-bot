package main

import (
	"log"
	"net/http"
	"poe-item-bot/internal/itemAPI"
)

func main() {
	// poe := api.New()
	// ladder, err := poe.GetLadder("SSF Ritual HC")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Got ladder with %d entries!", ladder.Total)
	// character, err := poe.GetCharacterItems("Zizaran", "ZizaranSmarter")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, item := range character.Items {
	// 	searchString := "Movement Speed"
	// 	for _, mod := range item.ExplicitMods {
	// 		if strings.Contains(mod, searchString) {
	// 			fmt.Printf("%s -- %+v\n", item.Type, item.ExplicitMods)
	// 		}
	// 	}
	// }
	i := itemAPI.NewItemAPI()
	log.Fatal(http.ListenAndServe("localhost:8080", i))
}
