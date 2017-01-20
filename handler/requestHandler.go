package handler

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
	"database/sql"
	"Retail/priceManager/model"
	"Retail/priceManager/database"
	"Retail/priceManager/status"
	"log"
)

func SaveUpdatePrice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		name := c.Query("name")
		id := parseId(c)


		price := parsePrice(c)

		database.ValidateEntryPresence(db, id);

		fmt.Printf("updating %s to cost %.2f\n", name, price)

		updatedPrice := model.Price{}
		updatedPrice.Product_id = id
		updatedPrice.Product_name = name
		updatedPrice.Cost = price
		updatedPrice.Status = status.PENDING

		database.SavePriceInUpdateTable(db, &updatedPrice);
	}
}

func parsePrice(c *gin.Context) float64 {
	cost := c.Query("cost")

	price, err := strconv.ParseFloat(cost, 64)

	if (err != nil) {
		panic(fmt.Sprintf("Given cost %s is wrong", cost))
	}
	return price
}

func parseId(c *gin.Context) int {
	id, err := strconv.Atoi(c.Query("id"))

	if (err != nil) {
		log.Fatal("Id format wrong")
	}

	return id;
}

