package pantry_api

import (
	"cooking/m/v2/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	db := c.MustGet("db").(*gorm.DB)
	var produce database.Produce
	db.Find(&produce, id)

	producesList := make([]map[string]interface{}, 1)
	producesList[0] = produce.ToMap()

	c.JSON(http.StatusOK, gin.H{"produces": producesList})
}
