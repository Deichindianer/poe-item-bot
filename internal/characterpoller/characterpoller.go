package characterpoller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/Deichindianer/poe-item-bot/pkg/api"
)

type CharacterPoller struct {
	PollList   []*PollCharacter
	Characters []*CharacterWindow
	ticker     *time.Ticker
	*api.Client
}

type PollCharacter struct {
	AccountName   string
	CharacterName string
}

type CharacterWindow struct {
	Items       []Item    `json:"items"`
	Character   Character `json:"character"`
	AccountName string
}

type Socket struct {
	GroupID   int    `json:"group"`
	Attribute string `json:"attr"`
}

type ItemProperty struct {
	Name        string        `json:"name"`
	Values      []interface{} `json:"values"`
	DisplayMode int           `json:"displayMode"`
}

type FrameType int

type Item struct {
	// Names for some items may include markup. For example: <<set:MS>><<set:M>><<set:S>>Roth's Reach
	Name string `json:"name"`
	Type string `json:"typeLine"`

	Properties   []ItemProperty `json:"properties"`
	Requirements []ItemProperty `json:"requirements"`

	Sockets []Socket `json:"sockets"`

	ExplicitMods []string `json:"explicitMods"`
	ImplicitMods []string `json:"implicitMods"`
	UtilityMods  []string `json:"utilityMods"`
	EnchantMods  []string `json:"enchantMods"`
	CraftedMods  []string `json:"craftedMods"`
	CosmeticMods []string `json:"cosmeticMods"`

	ID          string    `json:"id"`
	FrameType   FrameType `json:"frameType"`
	InventoryID string    `json:"inventoryId"`
	// Maybe I'll come back to needing it but not at the moment
	// SocketedItems []Item    `json:"socketedItems"`
}

type Character struct {
	Name       string `json:"name"`
	League     string `json:"league"`
	ClassID    int    `json:"classId"`
	Ascendancy string `json:"ascendancyClass"`
	Class      string `json:"class"`
	Level      int    `json:"level"`
	Experience int64  `json:"experience"`
	LastActive bool   `json:"lastActive"`
}

type refreshResult struct {
	characterWindow *CharacterWindow
	err             error
}

func NewCharacterPoller(pollList []*PollCharacter) *CharacterPoller {
	client := api.New()
	client.Scheme = "http"
	client.Host = "www.pathofexile.com"
	return &CharacterPoller{
		PollList: pollList,
		Client:   client,
	}
}

func (c *CharacterPoller) Poll(duration time.Duration) {
	if duration < time.Minute {
		log.Printf("Reset poll duration from %s to 1 minute\n", duration)
		duration = time.Minute
	}
	c.ticker = time.NewTicker(duration)
	if err := c.refreshAllCharacterItems(); err != nil {
		log.Printf("got error: %s", err)
		log.Println(err)
	}
	go func() {
		for range c.ticker.C {
			if err := c.refreshAllCharacterItems(); err != nil {
				log.Printf("got error: %s", err)
				log.Println(err)
			}
		}
	}()
}

func (c *CharacterPoller) StopPoll() {
	c.ticker.Stop()
}

func (c *CharacterPoller) refreshAllCharacterItems() error {
	start := time.Now()
	var wg sync.WaitGroup
	charactersChan := make(chan *refreshResult, len(c.PollList))
	for _, pollItem := range c.PollList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, pollItem *PollCharacter) {
			c.getCharacterWindow(charactersChan, pollItem.AccountName, pollItem.CharacterName)
			log.Println("Got the characterWindow")
			wg.Done()
		}(&wg, pollItem)
	}
	log.Println("Waiting for refreshCharacterWindow routines to finish")
	wg.Wait()
	log.Println("refreshCharacterWindow finished!")
	close(charactersChan)
	// for result := range charactersChan {
	// 	if result.err != nil {
	// 		log.Println(result.err)
	// 		continue
	// 	}
	// 	c.Characters = append(c.Characters, result.characterWindow)
	// }
	fmt.Printf("Duration of RefreshAllCharacterItems: %s\n", time.Since(start))
	return nil
}

func (c *CharacterPoller) getCharacterWindow(rc chan<- *refreshResult, accountName, characterName string) {
	query := url.Values{}
	query.Set("character", characterName)
	query.Set("accountName", accountName)
	response, err := c.CallAPI("character-window/get-items", query.Encode())
	if err != nil {
		log.Printf("got error: %s", err)
		rc <- &refreshResult{nil, err}
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("got error: %s", err)
		rc <- &refreshResult{nil, err}
	}
	var characterWindow CharacterWindow
	err = json.Unmarshal(responseBytes, &characterWindow)
	if err != nil {
		log.Printf("got error: %s", err)
		rc <- &refreshResult{nil, err}
	}
	characterWindow.AccountName = accountName
	log.Printf("char window successfully pulled, seding to channel")
	rc <- &refreshResult{&characterWindow, nil}
}
