package routes

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"Retail/priceManager/handler"
)

func HandleRequest(db *sql.DB) {
	router := gin.Default()

	router.POST("/price/update", handler.SaveUpdatePrice(db))

	router.Run(":4000")
}

