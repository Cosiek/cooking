package pantry_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cooking/m/v2/pantry"
)

func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}

func StartServer(db *gorm.DB) {
	r := gin.Default()

	r.Use(DBMiddleware(db))

	r.GET("/produce", func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		var products []pantry.Product
		db.Find(&products)
		names := ""
		for _, pr := range products {
			names += pr.Name
		}
		c.JSON(http.StatusOK, gin.H{"produces": names})
	})
	})
	r.Run(":8081")
}
