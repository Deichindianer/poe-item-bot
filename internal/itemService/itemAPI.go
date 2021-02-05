package itemAPI

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Deichindianer/poe-item-bot/pkg/api"
)

type ItemAPI struct {
	mux *gin.Engine
	poe *api.Client
}

type SearchResult struct {
	Metadata interface{}
	Items    []struct {
		Type         string
		ExplicitMods []string
	}
}

func NewItemAPI() (*ItemAPI, error) {
	i := new(ItemAPI)
	i.poe = api.New()
	i.mux = gin.New()
	i.mux.GET("/search", i.search)
	return i, nil
}

func (i *ItemAPI) search(c *gin.Context) {
	modSearch := c.Query("mod")
	typeSearch := c.Query("type")
	result := SearchResult{}
	character, err := i.poe.GetCharacterItems("Zizaran", "ZizaranSmarter")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range character.Items {
		if item.InventoryID == typeSearch || typeSearch == "" {
			for idx, mod := range item.ExplicitMods {
				if strings.Contains(mod, modSearch) {
					var resultItem = []struct {
						Type         string
						ExplicitMods []string
					}{{
						Type:         item.Type,
						ExplicitMods: item.ExplicitMods,
					}}
					result.Items = append(result.Items, resultItem...)
					fmt.Printf("%s -- %+v\n", item.Type, item.ExplicitMods[idx])
				}
			}
		}
	}
	c.JSON(http.StatusOK, result)
	return
}

func (i *ItemAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.mux.ServeHTTP(w, r)
}
