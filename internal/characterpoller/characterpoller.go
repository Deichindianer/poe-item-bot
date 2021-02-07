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
	PollList []PollCharacter
	ticker   *time.Ticker
	*api.Client

	mut        sync.RWMutex
	characters []CharacterWindow
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
	Ascendancy int    `json:"ascendancyClass"`
	Class      string `json:"class"`
	Level      int    `json:"level"`
	Experience int64  `json:"experience"`
	LastActive bool   `json:"lastActive"`
}

type refreshResult struct {
	characterWindow CharacterWindow
	err             error
}

func NewCharacterPoller(pollList []PollCharacter) *CharacterPoller {
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
	c.refreshAllCharacterItems()
	go func() {
		for range c.ticker.C {
			c.refreshAllCharacterItems()
		}
	}()
}

func (c *CharacterPoller) StopPoll() {
	c.ticker.Stop()
}

func (c *CharacterPoller) GetCharacters() []CharacterWindow {
	chars := make([]CharacterWindow, 0)
	c.mut.RLock()
	defer c.mut.RUnlock()
	chars = append(chars, c.characters...)
	return chars
}

func (c *CharacterPoller) refreshAllCharacterItems() {
	start := time.Now()
	// var tmpCharacter
	var wg sync.WaitGroup
	var tmpCharacters []CharacterWindow
	resultChannel := make(chan *refreshResult)
	go func(resultChannel <-chan *refreshResult) {
		for result := range resultChannel {
			if result.err != nil {
				log.Println(result.err)
				continue
			}
			tmpCharacters = append(tmpCharacters, result.characterWindow)
		}
	}(resultChannel)
	for _, pollItem := range c.PollList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, pollItem PollCharacter) {
			c.getCharacterWindow(resultChannel, pollItem.AccountName, pollItem.CharacterName)
			wg.Done()
		}(&wg, pollItem)
	}
	wg.Wait()
	close(resultChannel)
	c.mut.Lock()
	defer c.mut.Unlock()
	c.characters = tmpCharacters
	fmt.Printf("Duration of RefreshAllCharacterItems: %s\n", time.Since(start))
}

func (c *CharacterPoller) getCharacterWindow(rc chan<- *refreshResult, accountName, characterName string) {
	var characterWindow CharacterWindow
	query := url.Values{}
	query.Set("character", characterName)
	query.Set("accountName", accountName)
	response, err := c.CallAPI("character-window/get-items", query.Encode())
	if err != nil {
		log.Printf("failed to call API: %s", err)
		rc <- &refreshResult{characterWindow, err}
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read body: %s", err)
		rc <- &refreshResult{characterWindow, err}
	}

	err = json.Unmarshal(responseBytes, &characterWindow)
	if err != nil {
		log.Printf("failed to unmarshal characterWindow: %s", err)
		rc <- &refreshResult{characterWindow, err}
	}
	characterWindow.AccountName = accountName
	log.Printf("char window successfully pulled, seding to channel")
	rc <- &refreshResult{characterWindow, nil}
}
