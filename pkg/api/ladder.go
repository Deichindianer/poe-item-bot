package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Ladder struct {
	Total       int64     `json:"total"`
	CachedSince time.Time `json:"cached_since"`
	Entries     []Entry   `json:"entries"`
}

type Entry struct {
	Rank      int             `json:"rank"`
	Dead      bool            `json:"dead"`
	Online    bool            `json:"online"`
	Public    bool            `json:"public"`
	Character LadderCharacter `json:"character"`
	Account   Account         `json:"account"`
}

type LadderCharacter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Class      string `json:"class"`
	Score      int64  `json:"score"`
	Experience int64  `json:"experience"`
}

type Account struct {
	Name       string     `json:"name"`
	Realm      string     `json:"realm"`
	Challenges Challenges `json:"challenges"`
	Twitch     Twitch     `json:"twitch,omitempty"`
}

type Challenges struct {
	Total int `json:"total"`
}

type Twitch struct {
	Name string `json:"name"`
}

func (c Client) GetLadder(ladderName string) (*Ladder, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api.pathofexile.com/ladders/standard", nil)
	if err != nil {
		return nil, err
	}
	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var ladder Ladder
	if err = json.Unmarshal(responseBytes, &ladder); err != nil {
		return nil, err
	}
	return &ladder, nil
}
