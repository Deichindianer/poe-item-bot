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
	ticker     *time.Ticker
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
		Client:     client,
	}
}

func (l *LadderPoller) Poll(duration time.Duration) {
	if duration < time.Minute {
		log.Printf("Reset poll duration from %s to 1 minute\n", duration)
		duration = time.Minute
	}
	l.ticker = time.NewTicker(duration)
	log.Println("Refreshing ladder...")
	if err := l.refreshLadder(); err != nil {
		log.Printf("%s", err)
	}
	go func() {
		for range l.ticker.C {
			if err := l.refreshLadder(); err != nil {
				log.Printf("%s", err)
			}
		}
	}()
}

func (l *LadderPoller) StopPoll() {
	l.ticker.Stop()
}

func (l *LadderPoller) refreshLadder() error {
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
