package itemservice

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Deichindianer/poe-item-bot/internal/characterpoller"
	"github.com/gin-gonic/gin"
)

func (i *ItemService) search(c *gin.Context) {
	typeSearchString := c.Query("type")
	modSearchString := c.Query("mod")
	// linkSearchString := c.Query("links")
	result := SearchResult{}
	characterList := i.characterPoller.GetCharacters()

	if typeSearchString != "" {
		for _, cw := range characterList {
			tsr := typeSearch(typeSearchString, cw.Items)
			if tsr != nil {
				var match Match
				match.AccountName = cw.AccountName
				match.Items = tsr
				result.Result = append(result.Result, match)
			}
		}
	} else if modSearchString != "" {
		for _, cw := range characterList {
			msr := modSearch(modSearchString, cw.Items)
			if msr != nil {
				var match Match
				match.AccountName = cw.AccountName
				match.Items = msr
				result.Result = append(result.Result, match)
			}
		}
	}
	// else if linkSearchString != "" {
	// 	for _, cw := range characterList {
	// 		lsr := modSearch(modSearchString, cw.Items)
	// 		if lsr != nil {
	// 			var match Match
	// 			match.AccountName = cw.AccountName
	// 			match.Items = lsr
	// 			result.Result = append(result.Result, match)
	// 		}
	// 	}
	// }
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

func linkSearch(search string, items []characterpoller.Item) []characterpoller.Item {
	var result []characterpoller.Item
	for _, item := range items {
		for _, link := range item.Sockets {
			if strconv.Itoa(link.GroupID) == search {
				result = append(result, item)
			}
		}
	}
	return result
}
