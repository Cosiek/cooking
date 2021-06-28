package pantry_api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	r.POST("/produce", createProduce)
	r.GET("/produce", readAllProduces)
	r.GET("/produce/:id", readOneProduce)
	r.PUT("/produce/:id", updateProduce)
	r.DELETE("/produce/:id", deleteProduce)

	r.Run(":8081")
}
