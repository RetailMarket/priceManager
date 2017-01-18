package routes

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
)

func HandleRequest() {
	router := gin.Default()

	router.GET("/price/update", saveUpdatedPrice)

	router.Run(":4000")
}

func saveUpdatedPrice(c *gin.Context) {
	name := c.Query("name")
	price := parsePrice(c)

	// can store price to price update table
	fmt.Printf("updating %s price to %.2f\n", name, price)
}

func parsePrice(c *gin.Context) float64 {
	price, err := strconv.ParseFloat(c.Query("price"), 64)

	if (err != nil) {
		panic(fmt.Sprintf("Given price %s is wrong", c.Query("price")))
	}
	return price
}