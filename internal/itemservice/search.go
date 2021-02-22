package itemservice

import (
	"net/http"
	"strings"

	"github.com/Deichindianer/poe-item-bot/internal/characterpoller"
	"github.com/gin-gonic/gin"
)

func (i *ItemService) search(c *gin.Context) {
	typeSearchString := c.Query("type")
	modSearchString := c.Query("mod")
	result := SearchResult{}
	characterList := i.characterPoller.GetCharacters()

	if typeSearchString != "" {
		for _, cw := range characterList {
			var match Match
			match.AccountName = cw.AccountName
			tsr := typeSearch(typeSearchString, cw.Items)
			match.Items = tsr
			result.Result = append(result.Result, match)
		}
	} else {
		for _, cw := range characterList {
			var match Match
			match.AccountName = cw.AccountName
			msr := modSearch(modSearchString, cw.Items)
			match.Items = msr
			result.Result = append(result.Result, match)
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
