package routes

import (
	"github.com/gin-gonic/gin"
	"Retail/priceManager/handler"
)

func HandleRequest() {
	router := gin.Default()

	router.POST("/price/update", handler.SaveUpdatePrice)

	router.Run(":4000")
}

