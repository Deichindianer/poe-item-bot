package itemservice

import (
	"errors"
	"net/http"
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
	Result   []Match
}

type Match struct {
	AccountName string
	Items       []interface{}
}

func NewItemService(ladderName string, limit, offset int) *ItemService {
	i := new(ItemService)
	i.characterPoller = characterpoller.NewCharacterPoller(nil)
	i.ladderPoller = ladderpoller.NewLadderPoller(ladderName, limit, offset)
	gin.SetMode(gin.ReleaseMode)
	i.mux = gin.New()
	i.mux.GET("/search", i.search)
	i.mux.GET("/newsearch", i.newSearch)
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

func (i *ItemService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.mux.ServeHTTP(w, r)
}
