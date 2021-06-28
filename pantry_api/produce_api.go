package pantry_api

import (
	"cooking/m/v2/database"
	"encoding/json"
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
Creates a new produce entry in db.

If data passed to create the product is invalid, then the response will contain a dictionary of error messages.
*/
func createProduce(c *gin.Context) {
	// get data from request
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request - " + err.Error()})
	}

	// parse json
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request - " + err.Error()})
	}

	// validate input
	produce, errors := database.GetProduceFromMap(result)
	if len(errors) == 0 {
		// save to db
		db := c.MustGet("db").(*gorm.DB)
		db.Save(produce)

		c.JSON(http.StatusOK, gin.H{"message": "Ok", "produce": produce.ToMap()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Fial", "errors": errors})
	}
}

/*
Returns a list of all produces in the database.

TODO: limit to only the producess owned by current user.
*/
func readAllProduces(c *gin.Context) {
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
func readOneProduce(c *gin.Context) {
	id := c.Param("id")

	produce := getProduceOr404(id, c)
	if produce == nil {
		return
	}

	producesList := make([]map[string]interface{}, 1)
	producesList[0] = produce.ToMap()

	c.JSON(http.StatusOK, gin.H{"produces": producesList})
}

/*
Updates a produce in db.

If data passed to create the product is invalid, then the response will contain a dictionary of error messages.
*/
func updateProduce(c *gin.Context) {
	// get data from request
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request - " + err.Error()})
	}

	// get produce from db
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	produce := getProduceOr404(id, c)
	if produce == nil {
		return
	}

	// parse json
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request - " + err.Error()})
	}

	// attach produce id
	result["id"] = produce.ID

	// validate input
	produce, errors := database.GetProduceFromMap(result)
	if len(errors) == 0 {
		// save to db
		db.Save(produce)

		c.JSON(http.StatusOK, gin.H{"message": "Ok", "produce": produce.ToMap()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Fial", "errors": errors})
	}
}

/*
Deletes the produce refered by id in the url
*/
func deleteProduce(c *gin.Context) {
	id := c.Param("id")

	produce := getProduceOr404(id, c)
	if produce == nil {
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Delete(&produce)

	c.JSON(http.StatusOK, gin.H{"message": `Produce {id} deleted.`})
}
