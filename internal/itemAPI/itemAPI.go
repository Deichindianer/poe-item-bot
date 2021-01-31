package itemAPI

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemAPI struct {
	mux *gin.Engine
}

type SearchResult struct{}

func NewItemAPI() (*ItemAPI, error) {
	i := new(ItemAPI)

	i.mux = gin.New()
	i.mux.GET("/search", i.search)
	return i, nil
}

func (i *ItemAPI) search(c *gin.Context) {
	searchString := c.Query("searchString")
	c.JSON(http.StatusOK, struct{ searchString string }{searchString})
}

func (i *ItemAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.mux.ServeHTTP(w, r)
}
