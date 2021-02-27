package itemservice

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Deichindianer/poe-item-bot/internal/characterpoller"
	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
)

func (i *ItemService) newSearch(c *gin.Context) {
	result := SearchResult{}
	searches := c.Request.URL.Query()
	for _, cw := range i.characterPoller.GetCharacters() {
		var match Match
		var sr []mapset.Set
		for searchType, searchValue := range searches {
			switch searchType {
			case "type":
				if len(searchValue) != 1 {
					c.String(http.StatusBadRequest, "Only one type search param is allowed.")
					return
				}
				tsrSlice := typeSearch(searchValue[0], cw.Items)
				tsr := mapset.NewSetFromSlice(tsrSlice)
				sr = append(sr, tsr)
			case "mod":
				msrSlice := newModSearch(searchValue, cw.Items)
				msr := mapset.NewSetFromSlice(msrSlice)
				sr = append(sr, msr)
			case "links":
				if len(searchValue) != 1 {
					c.String(http.StatusBadRequest, "Only one link search param is allowed.")
					return
				}
				lsrSlice := linkSearch(searchValue[0], cw.Items)
				lsr := mapset.NewSetFromSlice(lsrSlice)
				sr = append(sr, lsr)
			}
		}
		match.AccountName = cw.AccountName
		var searchResult mapset.Set
		for _, s := range sr {
			if s != nil {
				if searchResult != nil {
					searchResult = searchResult.Intersect(s)
				} else {
					searchResult = s
				}
			}
		}
		match.Items = searchResult.ToSlice()
		result.Result = append(result.Result, match)
	}
	return
}

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
		modSearchStringLower := strings.ToLower(modSearchString)
		for _, cw := range characterList {
			msr := modSearch(modSearchStringLower, cw.Items)
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
	c.String(http.StatusOK, result)
	return
}

func typeSearch(search string, items []characterpoller.Item) []interface{} {
	var result []interface{}
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.InventoryID), strings.ToLower(search)) {
			result = append(result, item)
		}
	}
	return result
}

func newModSearch(search []string, items []characterpoller.Item) []interface{} {
	var result []interface{}
	for _, item := range items {
		for _, mod := range item.ExplicitMods {
			isCompleteMatch := true
			for _, s := range search {
				if strings.Contains(strings.ToLower(mod), strings.ToLower(s)) {
					continue
				} else {
					isCompleteMatch = false
				}
			}
			if isCompleteMatch {
				result = append(result, item)
			}
		}
	}
	return result
}

func modSearch(search string, items []characterpoller.Item) []interface{} {
	var result []interface{}
	for _, item := range items {
		for _, mod := range item.ExplicitMods {
			if strings.Contains(mod, search) {
				result = append(result, item)
			}
		}
	}
	return result
}

func linkSearch(search string, items []characterpoller.Item) []interface{} {
	var result []interface{}
	for _, item := range items {
		for _, link := range item.Sockets {
			if strconv.Itoa(link.GroupID) == search {
				result = append(result, item)
			}
		}
	}
	return result
}
