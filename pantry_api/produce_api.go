package pantry_api

import (
	"cooking/m/v2/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
Fetches a produce from database by ID. Sets 404 response if no matching produce was found.
*/
func getProduceOr404(produceId string, ctx *gin.Context) *database.Produce {
	db := ctx.MustGet("db").(*gorm.DB)
	var produce database.Produce
	db.Find(&produce, produceId)

	if produce.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Produce not found."})
		return nil
	}
	return &produce
}

/*
Returns a list of all produces in the database.

TODO: limit to only the producess owned by current user.
*/
func readAll(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var produces []database.Produce
	db.Find(&produces)

	producesList := make([]map[string]interface{}, len(produces))
	for idx, produce := range produces {
		producesList[idx] = produce.ToMap()
	}
	c.JSON(http.StatusOK, gin.H{"produces": producesList})
}

/*
Returns a list with a single element - the produce refered by id in the url.

TODO: limit to only the producess owned by current user.
*/
func readOne(c *gin.Context) {
	id := c.Param("id")

	produce := getProduceOr404(id, c)
	if produce == nil {
		return
	}

	producesList := make([]map[string]interface{}, 1)
	producesList[0] = produce.ToMap()

	c.JSON(http.StatusOK, gin.H{"produces": producesList})
}

func delete(c *gin.Context) {
	id := c.Param("id")

	produce := getProduceOr404(id, c)
	if produce == nil {
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Delete(&produce)

	c.JSON(http.StatusOK, gin.H{"message": `Produce {id} deleted.`})
}
