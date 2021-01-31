package ladderpoller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Deichindianer/poe-item-bot/pkg/api"
)

type LadderPoller struct {
	Ladder     Ladder
	LadderName string
	StopChan   chan bool
	Updated    bool
	*api.Client
}

type Ladder struct {
	Total       int64     `json:"total"`
	CachedSince time.Time `json:"cached_since"`
	Entries     []Entry   `json:"entries"`
}

type Entry struct {
	Rank      int       `json:"rank"`
	Dead      bool      `json:"dead"`
	Online    bool      `json:"online"`
	Public    bool      `json:"public"`
	Character Character `json:"character"`
	Account   Account   `json:"account"`
}

type Character struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Class string `json:"class"`
}

type Account struct {
	Name  string `json:"name"`
	Realm string `json:"realm"`
}

func NewLadderPoller(ladderName string) *LadderPoller {
	client := api.New()
	client.Scheme = "http"
	client.Host = "api.pathofexile.com"
	return &LadderPoller{
		Ladder:     Ladder{},
		LadderName: ladderName,
		StopChan:   make(chan bool),
		Client:     client,
	}
}

func (l LadderPoller) Poll() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			if err := l.refreshLadder(); err != nil {
				log.Fatal(err)
			}
			l.Updated = true
			select {
			case <-ticker.C:
				continue
			case <-l.StopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (l LadderPoller) StopPoll() {
	close(l.StopChan)
}

func (l LadderPoller) refreshLadder() error {
	response, err := l.CallAPI(fmt.Sprintf("ladders/%s", l.LadderName), "limit=45&offset=0")
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Got not okay status from api: %s", response.Status)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(responseBytes, &l.Ladder); err != nil {
		return err
	}
	return nil
}
