package itemservice

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Deichindianer/poe-item-bot/internal/characterpoller"
	"github.com/Deichindianer/poe-item-bot/internal/ladderpoller"
)

type ItemService struct {
	mux             *gin.Engine
	characterPoller *characterpoller.CharacterPoller
	ladderPoller    *ladderpoller.LadderPoller
}

type SearchResult struct {
	Metadata interface{}
	Items    []struct {
		Type         string
		ExplicitMods []string
	}
}

func NewItemService(ladderName string) *ItemService {
	i := new(ItemService)
	i.characterPoller = characterpoller.NewCharacterPoller(nil)
	i.ladderPoller = ladderpoller.NewLadderPoller(ladderName)
	i.mux = gin.New()
	i.mux.GET("/search", i.search)
	return i
}

func (i *ItemService) Init() error {
	i.ladderPoller.Poll(time.Minute)

	var pollList []characterpoller.PollCharacter
	var err error
	var j int

	for j < 5 {
		log.Println("Trying to create pollList")
		pollList, err = i.getPollListFromLadder()
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	i.characterPoller.PollList = pollList
	log.Printf("PollList: %+v\n", i.characterPoller.PollList)
	if len(i.characterPoller.PollList) == 0 {
		return errors.New("character poller did not get a poll list")
	}
	i.characterPoller.Poll(time.Minute)
	return nil
}

func (i *ItemService) getPollListFromLadder() ([]characterpoller.PollCharacter, error) {
	var pollList []characterpoller.PollCharacter
	if len(i.ladderPoller.Ladder.Entries) == 0 {
		return nil, errors.New("no entries in ladder, cannot create poll list")
	}
	for _, entry := range i.ladderPoller.Ladder.Entries {
		pollList = append(
			pollList,
			characterpoller.PollCharacter{
				AccountName:   entry.Account.Name,
				CharacterName: entry.Character.Name,
			},
		)
	}
	return pollList, nil
}

func (i *ItemService) search(c *gin.Context) {
	// modSearch := c.Query("mod")
	// typeSearch := c.Query("type")
	// result := SearchResult{}
	// character, err := i.poe.GetCharacterItems("Zizaran", "ZizaranSmarter")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, item := range character.Items {
	// 	if item.InventoryID == typeSearch || typeSearch == "" {
	// 		for idx, mod := range item.ExplicitMods {
	// 			if strings.Contains(mod, modSearch) {
	// 				var resultItem = []struct {
	// 					Type         string
	// 					ExplicitMods []string
	// 				}{{
	// 					Type:         item.Type,
	// 					ExplicitMods: item.ExplicitMods,
	// 				}}
	// 				result.Items = append(result.Items, resultItem...)
	// 				fmt.Printf("%s -- %+v\n", item.Type, item.ExplicitMods[idx])
	// 			}
	// 		}
	// 	}
	// }
	// c.JSON(http.StatusOK, result)
	result := fmt.Sprintf("Length of characters: %d", len(i.characterPoller.Characters))
	c.String(http.StatusOK, result)
	return
}

func (i *ItemService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.mux.ServeHTTP(w, r)
}
