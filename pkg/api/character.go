package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type CharacterWindow struct {
	Items     []Item    `json:"items"`
	Character Character `json:"character"`
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

	IsVerified             bool      `json:"verified"`
	Width                  int       `json:"w"`
	Height                 int       `json:"h"`
	ItemLevel              int       `json:"ilvl"`
	Icon                   string    `json:"icon"`
	League                 string    `json:"league"`
	ID                     string    `json:"id"`
	IsIdentified           bool      `json:"identified"`
	IsCorrupted            bool      `json:"corrupted"`
	IsLockedToCharacter    bool      `json:"lockedToCharacter"`
	IsSupport              bool      `json:"support"`
	DescriptionText        string    `json:"descrText"`
	SecondDescriptionText  string    `json:"secDescrText"`
	FlavorText             []string  `json:"flavourText"`
	ArtFilename            string    `json:"artFilename"`
	FrameType              FrameType `json:"frameType"`
	StackSize              int       `json:"stackSize"`
	MaxStackSize           int       `json:"maxStackSize"`
	X                      int       `json:"x"`
	Y                      int       `json:"y"`
	InventoryID            string    `json:"inventoryId"`
	SocketedItems          []Item    `json:"socketedItems"`
	IsRelic                bool      `json:"isRelic"`
	TalismanTier           int       `json:"talismanTier"`
	ProphecyText           string    `json:"prophecyText"`
	ProphecyDifficultyText string    `json:"prophecyDiffText"`
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

func (c *Client) GetCharacterItems(accountName string, characterName string) (*CharacterWindow, error) {
	characterQuery := fmt.Sprintf("character=%s&accountName=%s", characterName, accountName)
	response, err := c.CallAPI("character-window/get-items", characterQuery)
	if err != nil {
		return nil, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var character CharacterWindow
	err = json.Unmarshal(responseBytes, &character)
	return &character, nil
}
