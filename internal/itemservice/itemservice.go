package itemservice

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Deichindianer/poe-item-bot/internal/characterpoller"
	"github.com/Deichindianer/poe-item-bot/internal/ladderpoller"
)

type ItemService struct {
	mux             *gin.Engine
	characterCache  []characterpoller.Character
	characterPoller *characterpoller.CharacterPoller
	ladderPoller    *ladderpoller.LadderPoller
}

type SearchResult struct {
	Metadata interface{}
	Items    []characterpoller.Item
}

func NewItemService(ladderName string, limit, offset int) *ItemService {
	i := new(ItemService)
	i.characterPoller = characterpoller.NewCharacterPoller(nil)
	i.ladderPoller = ladderpoller.NewLadderPoller(ladderName, limit, offset)
	gin.SetMode(gin.ReleaseMode)
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
	log.Debugf("PollList: %+v\n", i.characterPoller.PollList)
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
	typeSearchString := c.Query("type")
	modSearchString := c.Query("mod")
	result := SearchResult{}
	characterList := i.characterPoller.GetCharacters()

	if typeSearchString != "" {
		for _, cw := range characterList {
			tsr := typeSearch(typeSearchString, cw.Items)
			result.Items = append(result.Items, tsr...)
		}
	} else {
		for _, cw := range characterList {
			msr := modSearch(modSearchString, cw.Items)
			result.Items = append(result.Items, msr...)
		}
	}
	c.JSON(http.StatusOK, result)
	return
}

func typeSearch(search string, items []characterpoller.Item) []characterpoller.Item {
	var result []characterpoller.Item
	for _, item := range items {
		if strings.ToLower(item.InventoryID) == strings.ToLower(search) {
			result = append(result, item)
		}
	}
	return result
}

func modSearch(search string, items []characterpoller.Item) []characterpoller.Item {
	var result []characterpoller.Item
	for _, item := range items {
		for _, mod := range item.ExplicitMods {
			if strings.Contains(strings.ToLower(mod), strings.ToLower(search)) {
				result = append(result, item)
			}
		}
	}
	return result
}

func (i *ItemService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.mux.ServeHTTP(w, r)
}
